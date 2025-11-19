package builder

import (
	"context"

	"github.com/TonyGLL/gofetch/internal/analysis"
	"github.com/TonyGLL/gofetch/internal/config"
	"github.com/TonyGLL/gofetch/internal/indexer"
	"github.com/TonyGLL/gofetch/pkg/storage"
)

// NewMongoStore creates a new MongoStore instance.
func NewMongoStore(ctx context.Context, cfg config.Config) (*storage.MongoStore, error) {
	return storage.NewMongoStore(ctx, cfg.MongoURI, cfg.DBName)
}

// NewAnalyzer creates a new Analyzer instance.
func NewAnalyzer() *analysis.Analyzer {
	return analysis.NewFromEnv()
}

// NewIndexer creates a new Indexer instance.
func NewIndexer(analyzer *analysis.Analyzer, store *storage.MongoStore) *indexer.Indexer {
	return indexer.NewIndexer(analyzer, store)
}
