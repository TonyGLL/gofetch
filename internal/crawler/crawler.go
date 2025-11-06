package crawler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
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
		host, err := extractHost(url)
		if err != nil {
			fmt.Printf("Failed to extract host from URL %s: %v\n", url, err)
			continue
		}
		robotsTxt, err := fetchRobotsTxt(host)
		if err != nil {
			fmt.Printf("Failed to fetch robots.txt for %s: %v\n", host, err)
			continue
		}
		fmt.Printf("robots.txt for %s:\n%s\n", host, string(robotsTxt))
		// Further crawling logic would go here
	}
}

func extractHost(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	host := parsed.Hostname()
	if host == "" {
		return "", fmt.Errorf("no host found in URL: %s", rawURL)
	}

	return host, nil
}

// fetchRobotsTxt fetches robots.txt safely, validating the host and handling errors properly.
func fetchRobotsTxt(host string) ([]byte, error) {
	// === 1. Validate and sanitize the host (prevent SSRF) ===
	parsedURL, err := url.Parse("https://" + host)
	if err != nil {
		return nil, fmt.Errorf("invalid host: %w", err)
	}
	if parsedURL.Hostname() != host {
		return nil, fmt.Errorf("invalid host: contains scheme or path")
	}
	if isPrivateOrLocal(host) {
		return nil, fmt.Errorf("host %s is private or local, blocked for security", host)
	}

	// === 2. Build a safe URL ===
	robotsURL := "https://" + host + "/robots.txt"

	// === 3. Create a context with timeout (RECOMMENDED) ===
	seconds := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), seconds)
	defer cancel() // Important!

	// === 4. Create the request WITH context ===
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, robotsURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "MyCrawler/1.0 (+https://example.com/bot)")

	// === 5. HTTP client (also respects the context) ===
	client := &http.Client{
		// Timeout ya no es necesario si usas contexto, pero puedes dejarlo
	}

	// === 6. Execute with Do ===
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch robots.txt: %w", err)
	}
	defer resp.Body.Close()

	// === 7. Read body ===
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read robots.txt body: %w", err)
	}

	return body, nil
}

func isPrivateOrLocal(host string) bool {
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}
