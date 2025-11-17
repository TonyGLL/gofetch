package main

import (
	"fmt"
	"time"

	"github.com/TonyGLL/gofetch/internal/crawler"
)

func main() {
	start := time.Now()

	// Application entry point
	depth := 1
	fmt.Println("Crawler application started")
	crawlerInst := crawler.NewCrawler([]string{
		"https://go.dev/",
	}, depth)

	crawlerInst.Crawl()

	elapsed := time.Since(start) // Calculate elapsed time
	fmt.Printf("Elapsed: %s\n", elapsed)
}
