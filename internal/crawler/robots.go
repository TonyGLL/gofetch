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
	myUA = strings.ToLower(strings.TrimSpace(myUA))
	rules := NewRobotRules(myUA)
	star := NewRobotRules("*")

	scanner := bufio.NewScanner(bytes.NewReader(data))
	var currentUA string
	inMyBlock := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		field, value, ok := parseFieldValue(line)
		if !ok {
			continue
		}

		switch field {
		case "user-agent":
			currentUA, inMyBlock = handleUserAgent(value, myUA)
		case "disallow":
			handleDirective(&rules, &star, inMyBlock, currentUA == "*", value, applyDisallow)
		case "allow":
			handleDirective(&rules, &star, inMyBlock, currentUA == "*", value, applyAllow)
		case "crawl-delay":
			handleCrawlDelay(&rules, &star, inMyBlock, currentUA == "*", value)
		}
	}

	return selectApplicableRules(rules, star)
}

// parseFieldValue extrae field:value de una línea
func parseFieldValue(line string) (field, value string, ok bool) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return strings.TrimSpace(strings.ToLower(parts[0])), strings.TrimSpace(parts[1]), true
}

// handleUserAgent actualiza el UA actual y si aplica al bot
func handleUserAgent(value, myUA string) (currentUA string, inMyBlock bool) {
	currentUA = strings.ToLower(value)
	return currentUA, strings.EqualFold(currentUA, myUA)
}

// applyDisallow procesa una línea Disallow
func applyDisallow(r *RobotRules, value string) {
	if value == "" {
		r.Disallows = nil
		return
	}
	if strings.HasPrefix(value, "/") {
		r.Disallows = append(r.Disallows, value)
	}
}

// applyAllow procesa una línea Allow
func applyAllow(r *RobotRules, value string) {
	if strings.HasPrefix(value, "/") {
		r.Allows = append(r.Allows, value)
	}
}

// handleDirective aplica una directiva (Allow/Disallow) al bloque correcto
type directiveFunc func(*RobotRules, string)

func handleDirective(
	target, fallback **RobotRules,
	inTarget, inFallback bool,
	value string,
	fn directiveFunc,
) {
	if inTarget {
		fn(*target, value)
	} else if inFallback {
		fn(*fallback, value)
	}
}

// handleCrawlDelay procesa crawl-delay
func handleCrawlDelay(target, fallback **RobotRules, inTarget, inFallback bool, value string) {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return
	}
	if inTarget {
		(*target).CrawlDelay = f
	} else if inFallback {
		(*fallback).CrawlDelay = f
	}
}

// selectApplicableRules decide qué reglas devolver
func selectApplicableRules(rules, star *RobotRules) *RobotRules {
	if len(rules.Allows) > 0 || len(rules.Disallows) > 0 || rules.CrawlDelay > 0 {
		rules.AppliesToMe = true
		return rules
	}
	star.AppliesToMe = true
	return star
}
