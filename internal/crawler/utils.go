package crawler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func fetchRobotsTxt(host, userAgent string) ([]byte, error) {
	robotsURL := fmt.Sprintf("https://%s/robots.txt", host)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", robotsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func resolveURL(base, href string) string {
	if strings.HasPrefix(href, "http") {
		return href
	}
	u, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}
	return baseURL.ResolveReference(u).String()
}
