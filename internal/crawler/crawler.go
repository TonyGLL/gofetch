package crawler

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Crawler struct {
	urls    []string
	visited map[string]bool
	depth   int
}

func NewCrawler(urls []string, depth int) *Crawler {
	return &Crawler{
		urls:    urls,
		visited: make(map[string]bool),
		depth:   depth,
	}
}

func (c *Crawler) Crawl() {
	for _, url := range c.urls {
		fmt.Printf("Crawling URL: %s\n", url)
		// Fetch and parse robots.txt
		host := extractHost(url)
		robotsTxt, err := fetchRobotsTxt(host)
		if err != nil {
			fmt.Printf("Failed to fetch robots.txt for %s: %v\n", host, err)
			continue
		}
		fmt.Printf("robots.txt for %s:\n%s\n", host, string(robotsTxt))
		// Further crawling logic would go here
	}
}

func extractHost(url string) string {
	// Simple extraction of host from URL
	// In production code, use url.Parse from net/url package
	if len(url) > 8 && url[:8] == "https://" {
		url = url[8:]
	} else if len(url) > 7 && url[:7] == "http://" {
		url = url[7:]
	}
	for i, ch := range url {
		if ch == '/' {
			return url[:i]
		}
	}
	return url
}

func fetchRobotsTxt(host string) ([]byte, error) {
	url := "https://" + host + "/robots.txt"
	resp, err := http.NewRequestWithContext(context.Background(), "GET", url, http.NoBody)
	if err != nil {
		return nil, err // o permitir todo si falla
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
