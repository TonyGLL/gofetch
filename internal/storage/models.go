package storage

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Document (no changes)
type Document struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	URL        string             `bson:"url"`
	Title      string             `bson:"title"`
	Content    string             `bson:"content"`
	IndexedAt  time.Time          `bson:"indexed_at"`
	ModifiedAt time.Time          `bson:"modified_at"`
	FilePath   string             `bson:"file_path"`
}

// Posting (with the Positions field added)
type Posting struct {
	DocID     primitive.ObjectID `bson:"doc_id"`
	Frequency int                `bson:"tf"`        // Changed to 'tf' by convention (term frequency)
	Positions []int              `bson:"positions"` // ADDED: for phrase searches
}

// InvertedIndexEntry (with the DF field added)
type InvertedIndexEntry struct {
	// We use the term as the _id for faster lookups and to ensure uniqueness.
	Term     string    `bson:"_id"`
	Postings []Posting `bson:"postings"`
	DF       int       `bson:"df"` // ADDED: Document Frequency
}

// IndexStats (no changes)
type IndexStats struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	TotalDocuments int64              `bson:"total_documents"`
	TotalTerms     int64              `bson:"total_terms"` // This field is more complex to compute; we'll focus on TotalDocuments
	LastIndexedAt  time.Time          `bson:"last_indexed_at"`
}
