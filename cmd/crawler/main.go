package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TonyGLL/gofetch/internal/builder"
	"github.com/TonyGLL/gofetch/internal/config"
	"github.com/TonyGLL/gofetch/internal/crawler"
)

func main() {
	start := time.Now()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	ctx := context.Background()
	store, err := builder.NewMongoStore(ctx, cfg)
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

	// Application entry point
	fmt.Println("Crawler application started")
	crawlerInst := crawler.NewCrawler(
		cfg.Crawler.URLs,
		cfg.Crawler.MaxDepth,
		idx,
		ctx,
	)

	crawlerInst.Crawl()

	elapsed := time.Since(start) // Calculate elapsed time
	fmt.Printf("Elapsed: %s\n", elapsed)
}
