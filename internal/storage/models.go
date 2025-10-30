package storage

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Document (sin cambios)
type Document struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	URL       string             `bson:"url"`
	Content   string             `bson:"content"`
	IndexedAt time.Time          `bson:"indexed_at"`
}

// Posting (con el campo Positions añadido)
type Posting struct {
	DocID     primitive.ObjectID `bson:"doc_id"`
	Frequency int                `bson:"tf"`        // Cambiado a 'tf' por convención (term frequency)
	Positions []int              `bson:"positions"` // AÑADIDO: para búsquedas de frases
}

// InvertedIndexEntry (con el campo DF añadido)
type InvertedIndexEntry struct {
	// Usamos el término como el _id para búsquedas más rápidas y para garantizar unicidad.
	Term     string    `bson:"_id"`
	Postings []Posting `bson:"postings"`
	DF       int       `bson:"df"` // AÑADIDO: Document Frequency
}

// IndexStats (sin cambios)
type IndexStats struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	TotalDocuments int64              `bson:"total_documents"`
	TotalTerms     int64              `bson:"total_terms"` // Este campo es más complejo de calcular, nos centraremos en TotalDocuments
	LastIndexedAt  time.Time          `bson:"last_indexed_at"`
}
