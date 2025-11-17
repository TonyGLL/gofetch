package analysis

import (
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/TonyGLL/gofetch/pkg/storage"
	"github.com/kljensen/snowball"
)

type Analyzer struct {
	stopwords map[string]struct{}
	language  string
}

func NewEnglishAnalyzer() *Analyzer {
	return New(storage.EnglishStopwords, "english")
}

func NewSpanishAnalyzer() *Analyzer {
	return New(storage.SpanishStopwords, "spanish")
}

func New(stopwords []string, language string) *Analyzer {
	return &Analyzer{
		stopwords: buildStopwordSet(stopwords),
		language:  language,
	}
}

func (a *Analyzer) Analyze(text string) []string {
	// Tokenize input text
	words := tokenize(text)

	// Normalize every word, filter stopwords, and stem
	result := []string{}
	for _, word := range words {
		normWord := normalize(word)
		if normWord == "" {
			continue
		}
		if _, isStopword := a.stopwords[normWord]; isStopword {
			continue
		}

		stemmedWord, err := snowball.Stem(normWord, a.language, true)
		if err != nil {
			// If stemming fails, use the normalized word as a fallback.
			result = append(result, normWord)
		} else {
			result = append(result, stemmedWord)
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

func NewFromEnv() *Analyzer {
	lang := os.Getenv("ANALYZER_LANGUAGE")
	switch strings.ToLower(lang) {
	case "spanish":
		return NewSpanishAnalyzer()
	default:
		return NewEnglishAnalyzer()
	}
}
