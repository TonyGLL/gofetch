package main

import (
	"fmt"

	"github.com/TonyGLL/gofetch/internal/crawler"
)

func main() {
	// Application entry point
	fmt.Println("Crawler application started")
	crawler.Crawl()
}
