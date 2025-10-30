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

// indexPayload es la estructura de datos que los workers envían al writer.
type indexPayload struct {
	Doc       storage.Document
	Freqs     map[string]int
	Positions map[string][]int
	FilePath  string
}

// Indexer encapsula la lógica de indexación.
type Indexer struct {
	analyzer    *analysis.Analyzer
	mongo_store *storage.MongoStore
}

// NewIndexer crea una nueva instancia del Indexer.
func NewIndexer(analyzer *analysis.Analyzer, mongo_store *storage.MongoStore) *Indexer {
	return &Indexer{
		analyzer:    analyzer,
		mongo_store: mongo_store,
	}
}

// IndexDirectory realiza el pipeline concurrente para indexar archivos de un directorio.
func (idx *Indexer) IndexDirectory(dirPath string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan string, 100)
	results := make(chan indexPayload, 100)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	workerCount := runtime.NumCPU()

	// 1. Iniciar Workers
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go idx.worker(ctx, &wg, jobs, results)
	}

	// 2. Iniciar Writer
	writeDone := make(chan struct{})
	go idx.writer(ctx, results, errCh, cancel, writeDone)

	// 3. Iniciar Producer
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

	// 4. Esperar y Sincronizar
	wg.Wait()
	close(results)

	select {
	case <-writeDone:
	case <-ctx.Done():
	}

	select {
	case err := <-errCh:
		return fmt.Errorf("indexing failed: %w", err)
	default:
		return nil
	}
}

// worker es la lógica que ejecuta cada goroutine del pool.
func (idx *Indexer) worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, results chan<- indexPayload) {
	defer wg.Done()
	for path := range jobs {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading file %s: %v\n", path, err)
				continue
			}
			text := string(data)
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

			payload := indexPayload{
				Doc: storage.Document{
					ID:        primitive.NewObjectID(),
					URL:       path,
					Content:   text,
					IndexedAt: time.Now(),
				},
				Freqs:     freqs,
				Positions: positions,
				FilePath:  path,
			}

			select {
			case results <- payload:
			case <-ctx.Done():
				return
			}
		}
	}
}

// writer consume los resultados y los escribe en MongoDB en lotes.
func (idx *Indexer) writer(ctx context.Context, results <-chan indexPayload, errCh chan<- error, cancel context.CancelFunc, done chan<- struct{}) {
	defer close(done)

	const BATCH_SIZE = 100
	const BATCH_TIMEOUT = 5 * time.Second
	batch := make([]indexPayload, 0, BATCH_SIZE)
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
			// Actualizar stats de forma incremental
			if err := idx.mongo_store.UpdateIndexStats(context.Background(), totalDocsInBatch); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to update index stats: %v\n", err)
			}
			batch = batch[:0] // Resetear el lote
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case payload, ok := <-results:
			if !ok {
				flushBatch() // Escribir el último lote
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

// writeBatch construye y ejecuta las operaciones BulkWrite para un lote de payloads.
func (idx *Indexer) writeBatch(ctx context.Context, batch []indexPayload) error {
	if len(batch) == 0 {
		return nil
	}

	docModels := make([]mongo.WriteModel, 0, len(batch))
	termModels := make([]mongo.WriteModel, 0, len(batch)*20)

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

// reportError envía un error al canal de errores sin bloquearse.
func reportError(errCh chan<- error, err error) {
	select {
	case errCh <- err:
	default:
	}
}
