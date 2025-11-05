package crawler

import (
	"fmt"
)

type RobotRules struct {
	UserAgent   string
	Disallows   []string
	Allows      []string
	CrawlDelay  float64
	AppliesToMe bool
}

func NewRobotRules(userAgent string) *RobotRules {
	return &RobotRules{
		UserAgent:   userAgent,
		Disallows:   []string{},
		Allows:      []string{},
		CrawlDelay:  0,
		AppliesToMe: false,
	}
}

func (r *RobotRules) IsAllowed(path string) bool {
	// Check disallows first
	for _, disallow := range r.Disallows {
		if matchPath(path, disallow) {
			// If there's a matching allow rule, it's allowed
			for _, allow := range r.Allows {
				if matchPath(path, allow) {
					return true
				}
			}
			return false
		}
	}
	return true
}

func matchPath(path, rule string) bool {
	// Simple prefix match for now
	return rule != "" && len(path) >= len(rule) && path[:len(rule)] == rule
}

func (r *RobotRules) String() string {
	result := "User-agent: " + r.UserAgent + "\n"
	for _, disallow := range r.Disallows {
		result += "Disallow: " + disallow + "\n"
	}
	for _, allow := range r.Allows {
		result += "Allow: " + allow + "\n"
	}
	if r.CrawlDelay > 0 {
		result += "Crawl-delay: " + fmt.Sprintf("%f", r.CrawlDelay) + "\n"
	}
	return result
}
