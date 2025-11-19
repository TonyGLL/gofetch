package main

import (
	"context"
	"fmt"
	"log"

	"github.com/TonyGLL/gofetch/internal/builder"
	"github.com/TonyGLL/gofetch/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	ctx := context.Background()
	store, err := builder.NewMongoStore(ctx, &cfg)
	if err != nil {
		log.Fatalf("Error creating MongoStore: %v", err)
	}
	defer func() {
		if err := store.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	an := builder.NewAnalyzer()
	idx := builder.NewIndexer(an, store)
	if err := idx.IndexDirectory(cfg.Indexer.Path); err != nil {
		fmt.Printf("Index error: %v\n", err)
	} else {
		fmt.Println("Indexing completed OK")
	}
}
