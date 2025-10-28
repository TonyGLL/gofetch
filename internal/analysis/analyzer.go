package analysis

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/TonyGLL/gofetch/pkg/storage"
)

type Analyzer struct {
	stopwords map[string]struct{}
}

func NewEnglishAnalyzer() *Analyzer {
	return New(storage.EnglishStopwords)
}

func NewSpanishAnalyzer() *Analyzer {
	return New(storage.SpanishStopwords)
}

func New(stopwords []string) *Analyzer {
	return &Analyzer{
		stopwords: buildStopwordSet(stopwords),
	}
}

func (a *Analyzer) Analyze(text string) []string {
	// Tokenize input text
	words := tokenize(text)

	// Normalize every word and filter stopwords
	result := []string{}
	for _, word := range words {
		normWord := normalize(word)
		if normWord == "" {
			continue
		}
		if _, isStopword := a.stopwords[normWord]; !isStopword {
			result = append(result, normWord)
		}
	}
	return result
}

func tokenize(text string) []string {
	re := regexp.MustCompile(`[[:alpha:]]+`)
	words := re.FindAllString(text, -1)
	return words
}

func normalize(word string) string {
	lowerWord := strings.ToLower(word)
	var builder strings.Builder
	for _, r := range lowerWord {
		// Keep letters and hyphen
		if unicode.IsLetter(r) || r == '-' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func buildStopwordSet(stopwords []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, word := range stopwords {
		set[word] = struct{}{}
	}
	return set
}
