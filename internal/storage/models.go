package storage

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Document representa la estructura de un documento web guardado en la colección 'documents'.
// Es el contenido original que ha sido rastreado e indexado.
type Document struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	URL       string             `bson:"url"`
	Content   string             `bson:"content"` // Podríamos omitir esto en producción para ahorrar espacio
	IndexedAt time.Time          `bson:"indexed_at"`
}

// Posting representa una única aparición de un término en un documento.
// No es una colección de nivel superior, sino un sub-documento dentro de InvertedIndexEntry.
type Posting struct {
	DocID     primitive.ObjectID `bson:"doc_id"`
	Frequency int                `bson:"frequency"` // La frecuencia del término (TF) en este documento
}

// InvertedIndexEntry representa una entrada en el índice invertido.
// Cada documento corresponde a un término único y contiene una lista de todos los
// documentos donde aparece ese término.
type InvertedIndexEntry struct {
	// Nota: No usamos un ObjectID aquí. El término en sí es la clave natural.
	// Se debe crear un índice único en MongoDB sobre este campo para un rendimiento óptimo.
	Term     string    `bson:"term"`
	Postings []Posting `bson:"postings"`
}

// IndexStats representa las estadísticas globales del índice.
// Normalmente, solo habrá un documento de este tipo en la colección 'stats'.
type IndexStats struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	TotalDocuments int64              `bson:"total_documents"`
	TotalTerms     int64              `bson:"total_terms"`
	LastIndexedAt  time.Time          `bson:"last_indexed_at"`
}
