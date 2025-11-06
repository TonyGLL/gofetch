package crawler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {
	startURLs   []string
	maxDepth    int
	visited     sync.Map
	queue       chan *CrawlTask
	rulesCache  sync.Map // map[string]*RobotRules
	lastRequest sync.Map // map[string]time.Time
	client      *http.Client
	userAgent   string
	wg          sync.WaitGroup
	workerCount int
	results     []CrawlResult
	resultsMu   sync.Mutex
}

type CrawlTask struct {
	URL   string
	Depth int
}

type CrawlResult struct {
	URL         string
	Title       string
	StatusCode  int
	Depth       int
	AllowedByRP bool
}

const defaultUserAgent = "MyCrawler/1.0 (+https://example.com/bot)"
const defaultWorkerCount = 5
const secondsInMillisecond = 1000
const millisecondsInSecond = 10

func NewCrawler(startURLs []string, maxDepth int) *Crawler {
	c := &Crawler{
		startURLs:   startURLs,
		maxDepth:    maxDepth,
		queue:       make(chan *CrawlTask, secondsInMillisecond),
		client:      &http.Client{Timeout: millisecondsInSecond * time.Second},
		userAgent:   defaultUserAgent,
		workerCount: defaultWorkerCount,
	}
	return c
}

func (c *Crawler) Crawl() {
	// 1. Enqueue initial URLs
	for _, u := range c.startURLs {
		c.enqueue(u, 0)
	}

	// 2. Start workers (NO wg.Add!)
	for range c.workerCount {
		go c.worker()
	}

	// 3. Close queue when there are no more tasks
	go func() {
		c.wg.Wait()
		close(c.queue)
	}()

	// 4. Wait for everything to finish
	c.wg.Wait()

	// 5. Print results
	c.printResults()
}

func (c *Crawler) worker() {
	for task := range c.queue {
		c.crawlTask(task)
	}
}

func (c *Crawler) enqueue(rawURL string, depth int) {
	if depth > c.maxDepth {
		return
	}
	if _, loaded := c.visited.LoadOrStore(rawURL, true); loaded {
		return
	}
	c.wg.Add(1)
	c.queue <- &CrawlTask{URL: rawURL, Depth: depth}
}

func (c *Crawler) crawlTask(task *CrawlTask) {
	defer c.wg.Done() // ← ALWAYS runs!

	u, err := url.Parse(task.URL)
	if err != nil {
		log.Printf("Invalid URL: %s", task.URL)
		return
	}

	host := u.Hostname()
	rules := c.getRobotRules(host)
	path := u.Path
	if u.RawQuery != "" {
		path += "?" + u.RawQuery
	}

	if !rules.IsAllowed(path) {
		c.addResult(CrawlResult{
			URL:         task.URL,
			Depth:       task.Depth,
			AllowedByRP: false,
		})
		log.Printf("[BLOCKED by robots.txt] %s", task.URL)
		return
	}

	c.respectCrawlDelay(host, rules.CrawlDelay)

	req, err := http.NewRequestWithContext(context.Background(), "GET", task.URL, http.NoBody)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", c.userAgent)
	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("HTTP error %s: %v", task.URL, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Printf("Parse error %s: %v", task.URL, err)
		return
	}

	title := doc.Find("title").First().Text()
	c.addResult(CrawlResult{
		URL:         task.URL,
		Title:       strings.TrimSpace(title),
		StatusCode:  resp.StatusCode,
		Depth:       task.Depth,
		AllowedByRP: true,
	})

	log.Printf("[OK] [%d] Depth %d: %s", resp.StatusCode, task.Depth, task.URL)

	if task.Depth < c.maxDepth {
		baseURL := task.URL
		doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if !exists {
				return
			}
			absURL := resolveURL(baseURL, href)
			if absURL == "" {
				return
			}
			abs, err := url.Parse(absURL)
			if err != nil || abs.Hostname() != host {
				return
			}
			c.enqueue(absURL, task.Depth+1)
		})
	}
}

func (c *Crawler) getRobotRules(host string) *RobotRules {
	if val, ok := c.rulesCache.Load(host); ok {
		return val.(*RobotRules)
	}

	body, err := fetchRobotsTxt(host, c.userAgent)
	if err != nil {
		log.Printf("robots.txt error for %s: %v → allowing all", host, err)
		rules := NewRobotRules(c.userAgent)
		rules.AppliesToMe = true
		c.rulesCache.Store(host, rules)
		return rules
	}

	rules := ParseRobotsTxt(body, c.userAgent)
	c.rulesCache.Store(host, rules)
	return rules
}

func (c *Crawler) respectCrawlDelay(host string, delay float64) {
	if delay <= 0 {
		return
	}
	last, _ := c.lastRequest.LoadOrStore(host, time.Time{})
	lastTime := last.(time.Time)
	wait := time.Duration(delay*secondsInMillisecond) * time.Millisecond
	sleep := wait - time.Since(lastTime)
	if sleep > 0 {
		time.Sleep(sleep)
	}
	c.lastRequest.Store(host, time.Now())
}

func (c *Crawler) addResult(r CrawlResult) {
	c.resultsMu.Lock()
	c.results = append(c.results, r)
	c.resultsMu.Unlock()
}

func (c *Crawler) printResults() {
	fmt.Println("\n=== CRAWL SUMMARY ===")
	for _, r := range c.results {
		status := "OK"
		if !r.AllowedByRP {
			status = "BLOCKED"
		} else if r.StatusCode >= http.StatusBadRequest {
			status = fmt.Sprintf("ERROR %d", r.StatusCode)
		}
		fmt.Printf("[%s] Depth %d: %s\n", status, r.Depth, r.URL)
		if r.Title != "" {
			fmt.Printf("    Title: %s\n", r.Title)
		}
	}
	fmt.Printf("Total pages processed: %d\n", len(c.results))
}
