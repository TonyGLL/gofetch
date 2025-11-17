package main

import (
	"context"
	"fmt"
	"time"

	"github.com/TonyGLL/gofetch/internal/analysis"
	"github.com/TonyGLL/gofetch/internal/crawler"
	"github.com/TonyGLL/gofetch/internal/indexer"
	"github.com/TonyGLL/gofetch/internal/storage"
)

func main() {
	start := time.Now()

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

	an := analysis.NewEnglishAnalyzer()
	idx := indexer.NewIndexer(an, store)

	// Application entry point
	depth := 1
	fmt.Println("Crawler application started")
	crawlerInst := crawler.NewCrawler([]string{
		"https://go.dev/",
	}, depth, idx, ctx)

	crawlerInst.Crawl()

	elapsed := time.Since(start) // Calculate elapsed time
	fmt.Printf("Elapsed: %s\n", elapsed)
}
