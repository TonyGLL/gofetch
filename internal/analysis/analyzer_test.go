package analysis

import (
	"reflect"
	"testing"

	"github.com/TonyGLL/gofetch/pkg/storage"
)

type testCase struct {
	name           string   // Descriptive name of the test case
	language       string   // The language for the analyzer
	stopwords      []string // The Analyzer configuration to test
	inputText      string   // The input text to analyze
	expectedTokens []string // The exact result we expect
}

func TestAnalyzer_Analyze(t *testing.T) {
	testCases := []testCase{
		{
			name:           "Spanish text with stemming",
			language:       "spanish",
			stopwords:      storage.SpanishStopwords,
			inputText:      "Este es un TEXTO de prueba, ¡genial!",
			expectedTokens: []string{"text", "prueb", "genial"},
		},
		{
			name:           "English text with stemming",
			language:       "english",
			stopwords:      []string{},
			inputText:      "A simple text with running words to test.",
			expectedTokens: []string{"a", "simpl", "text", "with", "run", "word", "to", "test"},
		},
		{
			name:           "Input with only stopwords",
			language:       "english",
			stopwords:      storage.EnglishStopwords,
			inputText:      "It is a she or he",
			expectedTokens: []string{},
		},
		{
			name:           "Input with only punctuation and numbers",
			language:       "english",
			stopwords:      storage.EnglishStopwords,
			inputText:      "123.45, -¡!@#$%^&*()_+",
			expectedTokens: []string{},
		},
		{
			name:           "English stemming with common variations",
			language:       "english",
			stopwords:      storage.EnglishStopwords,
			inputText:      "running runner runs",
			expectedTokens: []string{"run", "runner", "run"},
		},
		{
			name:           "Spanish stemming with common variations",
			language:       "spanish",
			stopwords:      storage.SpanishStopwords,
			inputText:      "corriendo corredores corren",
			expectedTokens: []string{"corr", "corredor", "corr"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			analyzer := New(tc.stopwords, tc.language)
			tokens := analyzer.Analyze(tc.inputText)

			if !reflect.DeepEqual(tokens, tc.expectedTokens) {
				t.Errorf("Expected tokens %v, but got %v", tc.expectedTokens, tokens)
			}
		})
	}
}
