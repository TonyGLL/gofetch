package analysis

import (
	"reflect"
	"testing"

	"github.com/TonyGLL/gofetch/pkg/storage"
)

type testCase struct {
	name           string   // Descriptive name of the test case
	stopwords      []string // The Analyzer configuration to test
	inputText      string   // The input text to analyze
	expectedTokens []string // The exact result we expect
}

func TestAnalyzer_Analyze(t *testing.T) {
	testCases := []testCase{
		{
			name:           "Spanish text with Spanish stopwords",
			stopwords:      storage.SpanishStopwords,
			inputText:      "Este es un TEXTO de prueba, ¡genial!",
			expectedTokens: []string{"texto", "prueba", "genial"},
		},
		{
			name:           "Text with no stopwords configured",
			stopwords:      []string{},
			inputText:      "A simple text to test.",
			expectedTokens: []string{"a", "simple", "text", "to", "test"},
		},
		{
			name:           "Input with only stopwords",
			stopwords:      storage.EnglishStopwords,
			inputText:      "It is a she or he",
			expectedTokens: []string{},
		},
		{
			name:           "Input with only punctuation and numbers",
			stopwords:      storage.EnglishStopwords,
			inputText:      "123.45, -¡!@#$%^&*()_+",
			expectedTokens: []string{},
		},
		{
			name:           "Input with hyphenated word",
			stopwords:      []string{},
			inputText:      "This is state-of-the-art",
			expectedTokens: []string{"this", "is", "state", "of", "the", "art"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			analyzer := New(tc.stopwords)
			tokens := analyzer.Analyze(tc.inputText)

			if !reflect.DeepEqual(tokens, tc.expectedTokens) {
				t.Errorf("Expected tokens %v, but got %v", tc.expectedTokens, tokens)
			}
		})
	}
}
