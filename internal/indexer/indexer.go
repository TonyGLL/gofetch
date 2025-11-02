package indexer

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/TonyGLL/gofetch/internal/analysis"
	"github.com/TonyGLL/gofetch/internal/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// indexPayload is the data structure that workers send to the writer.
type indexPayload struct {
	Doc       storage.Document
	Freqs     map[string]int
	Positions map[string][]int
	FilePath  string
}

// Indexer encapsulates the indexing logic.
type Indexer struct {
	analyzer    *analysis.Analyzer
	mongo_store *storage.MongoStore
}

// NewIndexer creates a new Indexer instance.
func NewIndexer(analyzer *analysis.Analyzer, mongo_store *storage.MongoStore) *Indexer {
	return &Indexer{
		analyzer:    analyzer,
		mongo_store: mongo_store,
	}
}

// IndexDirectory runs the concurrent pipeline to index files in a directory.
func (idx *Indexer) IndexDirectory(dirPath string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	n := 100
	jobs := make(chan string, n)
	results := make(chan *indexPayload, n)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	workerCount := runtime.NumCPU()

	indexedFilesCh := make(chan string, n)

	// 1. Start workers
	wg.Add(workerCount)
	for range workerCount {
		go idx.worker(ctx, &wg, jobs, results, indexedFilesCh)
	}

	// 2. Start writer
	writeDone := make(chan struct{})
	go idx.writer(ctx, results, errCh, cancel, writeDone)

	// 3. Start producer
	go func() {
		defer close(jobs)
		walkErr := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if ext != ".txt" && ext != ".md" {
				return nil
			}
			select {
			case jobs <- path:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
		if walkErr != nil {
			reportError(errCh, walkErr)
			cancel()
		}
	}()

	// 4. Wait and synchronize
	wg.Wait()
	close(results)
	close(indexedFilesCh)

	select {
	case <-writeDone:
	case <-ctx.Done():
	}
	var indexedFiles []string
	for filePath := range indexedFilesCh {
		indexedFiles = append(indexedFiles, filePath)
	}
	select {
	case err := <-errCh:
		return fmt.Errorf("indexing failed: %w", err)
	default:
		fmt.Printf("Successfully indexed %d files:\n", len(indexedFiles))
		for _, file := range indexedFiles {
			fmt.Printf("- %s\n", file)
		}
		return nil
	}
}

func (idx *Indexer) processFile(ctx context.Context, path string) (*indexPayload, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}
	// Get file modification time
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("error getting file info for %s: %w", path, err)
	}
	modifiedAt := fileInfo.ModTime()

	// Check if the document is already indexed and unchanged
	existingDoc, err := idx.mongo_store.GetDocumentByPath(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("error checking existing document for %s: %w", path, err)
	}

	if existingDoc != nil {
		if !modifiedAt.After(existingDoc.ModifiedAt) {
			fmt.Printf("Skipping unchanged file: %s\n", path)
			return nil, nil // nil, nil indicates skipped file
		}
		// If the file has been modified, remove old postings and delete the old document
		oldTerms := idx.analyzer.Analyze(existingDoc.Content)
		if err := idx.mongo_store.RemovePostingsForDocument(ctx, existingDoc.ID, oldTerms); err != nil {
			return nil, fmt.Errorf("error removing old postings for %s: %w", path, err)
		}

		if err := idx.mongo_store.DeleteDocument(ctx, existingDoc.ID); err != nil {
			return nil, fmt.Errorf("error deleting existing document for %s: %w", path, err)
		}
	}

	text := string(data)

	// Title extraction logic
	title := ""
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			title = trimmedLine
			break
		}
	}

	if title == "" {
		title = filepath.Base(path) // Fallback to filename
	}

	tokens := idx.analyzer.Analyze(text)
	freqs := make(map[string]int)
	positions := make(map[string][]int)
	for i, token := range tokens {
		if token == "" {
			continue
		}

		freqs[token]++
		positions[token] = append(positions[token], i)
	}

	payload := &indexPayload{
		Doc: storage.Document{
			ID:         primitive.NewObjectID(),
			URL:        path,
			Title:      title, // Set the extracted title
			Content:    text,
			IndexedAt:  time.Now(),
			ModifiedAt: modifiedAt,
			FilePath:   path,
		},
		Freqs:     freqs,
		Positions: positions,
		FilePath:  path,
	}

	return payload, nil
}

// worker is the logic executed by each goroutine in the pool.

func (idx *Indexer) worker(
	ctx context.Context,
	wg *sync.WaitGroup,
	jobs <-chan string,
	results chan<- *indexPayload,
	indexedFilesCh chan<- string,
) {
	defer wg.Done()
	for path := range jobs {
		select {
		case <-ctx.Done():
			return
		default:
			payload, err := idx.processFile(ctx, path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}

			if payload == nil {
				continue // File was skipped
			}

			select {
			case results <- payload:
				select {
				case indexedFilesCh <- path:
				case <-ctx.Done():
					return
				}

			case <-ctx.Done():
				return
			}
		}
	}
}

// writer consumes results and writes them to MongoDB in batches.
func (idx *Indexer) writer(
	ctx context.Context,
	results <-chan *indexPayload,
	errCh chan<- error,
	cancel context.CancelFunc,
	done chan<- struct{},
) {
	defer close(done)

	const BATCH_SIZE = 100
	const BATCH_TIMEOUT = 5 * time.Second
	batch := make([]*indexPayload, 0, BATCH_SIZE)
	ticker := time.NewTicker(BATCH_TIMEOUT)
	defer ticker.Stop()

	totalDocsInBatch := int64(0)

	flushBatch := func() {
		if len(batch) == 0 {
			return
		}
		if err := idx.writeBatch(ctx, batch); err != nil {
			reportError(errCh, err)
			cancel()
		} else {
			totalDocsInBatch = int64(len(batch))
			// Incrementally update stats
			if err := idx.mongo_store.UpdateIndexStats(context.Background(), totalDocsInBatch); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to update index stats: %v\n", err)
			}
			batch = batch[:0] // Reset the batch
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case payload, ok := <-results:
			if !ok {
				flushBatch() // Write the last batch
				return
			}
			batch = append(batch, payload)
			if len(batch) >= BATCH_SIZE {
				flushBatch()
				ticker.Reset(BATCH_TIMEOUT)
			}
		case <-ticker.C:
			flushBatch()
		}
	}
}

// writeBatch builds and executes BulkWrite operations for a batch of payloads.
func (idx *Indexer) writeBatch(ctx context.Context, batch []*indexPayload) error {
	if len(batch) == 0 {
		return nil
	}
	mul := 20
	docModels := make([]mongo.WriteModel, 0, len(batch))
	termModels := make([]mongo.WriteModel, 0, len(batch)*mul)

	for _, payload := range batch {
		docModels = append(docModels, mongo.NewInsertOneModel().SetDocument(payload.Doc))

		for term, freq := range payload.Freqs {
			posting := storage.Posting{
				DocID:     payload.Doc.ID,
				Frequency: freq,
				Positions: payload.Positions[term],
			}
			model := mongo.NewUpdateOneModel().
				SetFilter(bson.M{"_id": term}).
				SetUpdate(bson.M{
					"$push": bson.M{"postings": posting},
					"$inc":  bson.M{"df": 1},
				}).
				SetUpsert(true)
			termModels = append(termModels, model)
		}
	}

	if err := idx.mongo_store.BulkWriteDocuments(ctx, docModels); err != nil {
		return fmt.Errorf("failed to bulk write documents: %w", err)
	}
	if err := idx.mongo_store.BulkWriteInvertedIndex(ctx, termModels); err != nil {
		return fmt.Errorf("failed to bulk write inverted index: %w", err)
	}

	fmt.Printf("Successfully indexed batch of %d documents.\n", len(batch))
	return nil
}

// reportError sends an error to the error channel without blocking.
func reportError(errCh chan<- error, err error) {
	select {
	case errCh <- err:
	default:
	}
}
