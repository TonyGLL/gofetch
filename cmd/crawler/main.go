package main

import (
	"fmt"

	"github.com/TonyGLL/gofetch/internal/crawler"
)

func main() {
	// Application entry point
	depth := 4
	fmt.Println("Crawler application started")
	crawlerInst := crawler.NewCrawler([]string{
		"https://go.dev/",
	}, depth)

	crawlerInst.Crawl()
}
