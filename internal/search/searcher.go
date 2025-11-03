package search

import (
	"context"
	"sort"

	"github.com/TonyGLL/gofetch/internal/analysis"
	"github.com/TonyGLL/gofetch/internal/ranking"
	"github.com/TonyGLL/gofetch/internal/storage"
)

// Searcher defines the interface for searching documents.
type Searcher interface {
	Search(ctx context.Context, query string, pagination storage.GetDocumentsFilter) (SearchDocumentResponse, error)
}

// SearchResult represents a single search result.
type SearchDocumentResponse struct {
	Data  []SearchResult `json:"data,omitempty"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Total int            `json:"total"`
}
type SearchResult struct {
	DocID string  `json:"docID"`
	Title string  `json:"title"`
	URL   string  `json:"url"`
	Score float64 `json:"-"`
}

// searcherImpl is the concrete implementation of the Searcher interface.
type searcherImpl struct {
	analyzer *analysis.Analyzer
	store    *storage.MongoStore
}

// NewSearcher creates a new instance of the searcher.
func NewSearcher(analyzer *analysis.Analyzer, store *storage.MongoStore) Searcher {
	return &searcherImpl{
		analyzer: analyzer,
		store:    store,
	}
}

// Search performs a search for the given query.
func (s *searcherImpl) Search(ctx context.Context, query string, pagination storage.GetDocumentsFilter) (SearchDocumentResponse, error) {
	// 1. Analyze the query string.
	queryTerms := s.analyzer.Analyze(query)

	// 2. Fetch index data for the query terms from the store.
	postings, err := s.store.GetPostingsForTerms(ctx, queryTerms)
	if err != nil {
		return SearchDocumentResponse{
			Page:  int(pagination.Page),
			Limit: int(pagination.Limit),
		}, err
	}

	// 3. Fetch global index stats for scoring.
	stats, err := s.store.GetIndexStats(ctx)
	if err != nil {
		return SearchDocumentResponse{
			Page:  int(pagination.Page),
			Limit: int(pagination.Limit),
		}, err
	}

	// 4. Score the documents using the TF-IDF ranker.
	scorer := ranking.NewTFIDFScorer(*stats)
	docScores := scorer.Score(queryTerms, postings)

	// 5. Fetch document metadata for the top-scoring documents.
	docIDs := make([]string, 0, len(docScores))
	for id := range docScores {
		docIDs = append(docIDs, id)
	}
	documents, total, err := s.store.GetDocuments(ctx, docIDs, pagination)
	if err != nil {
		return SearchDocumentResponse{
			Page:  int(pagination.Page),
			Limit: int(pagination.Limit),
		}, err
	}

	// 6. Build the final search results.
	results := make([]SearchResult, 0, len(documents))
	for _, doc := range documents {
		results = append(results, SearchResult{
			DocID: doc.ID.Hex(),
			Title: doc.Title,
			URL:   doc.URL,
			Score: docScores[doc.ID.Hex()],
		})
	}

	// 7. Sort the results by score in descending order.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	response := SearchDocumentResponse{
		Data:  results,
		Page:  int(pagination.Page),
		Limit: int(pagination.Limit),
		Total: total,
	}

	return response, nil
}
