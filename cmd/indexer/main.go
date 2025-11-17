package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/TonyGLL/gofetch/internal/analysis"
	"github.com/TonyGLL/gofetch/internal/indexer"
	"github.com/TonyGLL/gofetch/internal/storage"
)

func main() {
	store := storage.NewMongoStore()
	ctx := context.Background()
	if err := store.Connect(ctx); err != nil {
		panic(err)
	}
	defer func() {
		if err := store.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	args := flag.String("path", "data/index", "Path to store the index data")
	flag.Parse()

	an := analysis.NewFromEnv()
	idx := indexer.NewIndexer(an, store)
	if err := idx.IndexDirectory(*args); err != nil {
		fmt.Printf("Index error: %v\n", err)
	} else {
		fmt.Println("Indexing completed OK")
	}
}
