package ranking

import (
	"math"

	"github.com/TonyGLL/gofetch/pkg/storage"
)

// TFIDFScorer calculates scores based on the TF-IDF algorithm.
type TFIDFScorer struct {
	TotalDocuments int64
}

// NewTFIDFScorer creates a new scorer.
func NewTFIDFScorer(stats storage.IndexStats) *TFIDFScorer {
	return &TFIDFScorer{
		TotalDocuments: stats.TotalDocuments,
	}
}

// Score calculates the TF-IDF score for a set of documents based on a query.
func (s *TFIDFScorer) Score(queryTerms []string, postings map[string]storage.InvertedIndexEntry) map[string]float64 {
	docScores := make(map[string]float64)

	for _, term := range queryTerms {
		entry, ok := postings[term]
		if !ok {
			continue // Term not in index
		}

		idf := s.calculateIDF(entry.DF)

		for _, post := range entry.Postings {
			tf := float64(post.Frequency)
			docID := post.DocID.Hex()
			docScores[docID] += tf * idf
		}
	}

	return docScores
}

// calculateIDF calculates the Inverse Document Frequency for a term.
func (s *TFIDFScorer) calculateIDF(docFrequency int) float64 {
	if docFrequency == 0 || s.TotalDocuments == 0 {
		return 0
	}
	// Using the smooth IDF formula to avoid division by zero
	return math.Log(1 + (float64(s.TotalDocuments) / float64(docFrequency)))
}
