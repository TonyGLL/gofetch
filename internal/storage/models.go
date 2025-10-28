package storage

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Document represents the structure of a web document saved in the 'documents' collection.
// It is the original content that has been crawled and indexed.
type Document struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	URL       string             `bson:"url"`
	Content   string             `bson:"content"` // We could omit this in production to save space
	IndexedAt time.Time          `bson:"indexed_at"`
}

// Posting represents a single occurrence of a term in a document.
// It is not a top-level collection, but a sub-document within InvertedIndexEntry.
type Posting struct {
	DocID     primitive.ObjectID `bson:"doc_id"`
	Frequency int                `bson:"frequency"` // The frequency of the term (TF) in this document
}

// InvertedIndexEntry represents an entry in the inverted index.
// Each document corresponds to a unique term and contains a list of all the
// documents where that term appears.
type InvertedIndexEntry struct {
	// Note: We do not use an ObjectID here. The term itself is the natural key.
	// A unique index should be created in MongoDB on this field for optimal performance.
	Term     string    `bson:"term"`
	Postings []Posting `bson:"postings"`
}

// IndexStats represents the global statistics of the index.
// Normally, there will only be one document of this type in the 'stats' collection.
type IndexStats struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	TotalDocuments int64              `bson:"total_documents"`
	TotalTerms     int64              `bson:"total_terms"`
	LastIndexedAt  time.Time          `bson:"last_indexed_at"`
}
