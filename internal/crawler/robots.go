package crawler

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
)

type RobotRules struct {
	UserAgent   string
	Disallows   []string
	Allows      []string
	CrawlDelay  float64
	AppliesToMe bool
}

func NewRobotRules(ua string) *RobotRules {
	return &RobotRules{
		UserAgent:   ua,
		Disallows:   []string{},
		Allows:      []string{},
		CrawlDelay:  0,
		AppliesToMe: false,
	}
}

func (r *RobotRules) IsAllowed(path string) bool {
	if len(r.Disallows) == 0 && len(r.Allows) == 0 {
		return true
	}

	var bestAllow, bestDisallow string

	for _, a := range r.Allows {
		if strings.HasPrefix(path, a) && len(a) > len(bestAllow) {
			bestAllow = a
		}
	}
	for _, d := range r.Disallows {
		if strings.HasPrefix(path, d) && len(d) > len(bestDisallow) {
			bestDisallow = d
		}
	}

	if bestAllow != "" && bestDisallow != "" {
		if len(bestAllow) >= len(bestDisallow) {
			return true
		}
	}
	if bestAllow != "" {
		return true
	}
	if bestDisallow != "" {
		return false
	}
	return true
}

func ParseRobotsTxt(data []byte, myUA string) *RobotRules {
	rules := NewRobotRules(myUA)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var currentUA string
	inMyBlock := false
	starRules := NewRobotRules("*")

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		field := strings.TrimSpace(strings.ToLower(parts[0]))
		value := strings.TrimSpace(parts[1])

		switch field {
		case "user-agent":
			currentUA = strings.ToLower(value)
			inMyBlock = (currentUA == strings.ToLower(myUA))
		case "disallow":
			if value == "" {
				if inMyBlock {
					rules.Disallows = nil
				} else if currentUA == "*" {
					starRules.Disallows = nil
				}
			} else if strings.HasPrefix(value, "/") {
				if inMyBlock {
					rules.Disallows = append(rules.Disallows, value)
				} else if currentUA == "*" {
					starRules.Disallows = append(starRules.Disallows, value)
				}
			}
		case "allow":
			if strings.HasPrefix(value, "/") {
				if inMyBlock {
					rules.Allows = append(rules.Allows, value)
				} else if currentUA == "*" {
					starRules.Allows = append(starRules.Allows, value)
				}
			}
		case "crawl-delay":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				if inMyBlock {
					rules.CrawlDelay = f
				} else if currentUA == "*" {
					starRules.CrawlDelay = f
				}
			}
		}
	}

	if len(rules.Disallows)+len(rules.Allows) > 0 || rules.CrawlDelay > 0 {
		rules.AppliesToMe = true
		return rules
	}

	starRules.AppliesToMe = true
	return starRules
}
