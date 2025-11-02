package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	client             *mongo.Client
	database           *mongo.Database
	documentCollection *mongo.Collection
	indexCollection    *mongo.Collection
	statsCollection    *mongo.Collection
}

// NewMongoStore is a constructor function that initializes an instance of MongoStore.
// The actual connection is established through the Connect() method.
func NewMongoStore() *MongoStore {
	return &MongoStore{}
}

// Connect establishes the connection with MongoDB using a URI from an environment variable.
// It also performs a ping to verify the connection and prepares the collection handlers.
func (s *MongoStore) Connect(ctx context.Context) error {
	// Detailed logic for:
	// 1. Read os.Getenv("MONGO_URI")
	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "gofetch"
	}
	// 2. Configure clientOptions
	clientOptions := options.Client().ApplyURI(mongoURI)
	// 3. Call mongo.Connect(ctx, clientOptions)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	// 4. Call client.Ping(ctx, nil) to verify
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}
	// 5. If everything goes well, initialize the struct fields:
	s.client = client
	s.database = client.Database(dbName)
	s.documentCollection = s.database.Collection("documents")
	s.indexCollection = s.database.Collection("inverted_index")
	s.statsCollection = s.database.Collection("stats")

	fmt.Println("Connected to MongoDB successfully.")
	return nil
}

// Disconnect safely closes the database connection.
func (s *MongoStore) Disconnect(ctx context.Context) error {
	if s.client == nil {
		return nil
	}
	fmt.Println("Disconnected from MongoDB.")
	return s.client.Disconnect(ctx)
}

// AddDocument inserts a new document into the 'documents' collection.
// Returns the ID of the inserted document as a hexadecimal string.
func (s *MongoStore) AddDocument(ctx context.Context, doc *Document) (string, error) {
	result, err := s.documentCollection.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// UpsertTerm updates or inserts an entry in the inverted index.
// It searches for a term and adds the posting to its list. If the term does not exist,
// it creates it along with the initial posting. This operation is atomic.
func (s *MongoStore) UpsertTerm(ctx context.Context, term string, posting Posting) error {
	// The filter identifies the document by its term, which is the _id.
	filter := bson.M{"_id": term}

	// The update operation does two things:
	// 1. $push: Adds the new posting to the 'postings' array.
	// 2. $inc: Increments the document frequency (df) by 1.
	update := bson.M{
		"$push": bson.M{"postings": posting},
		"$inc":  bson.M{"df": 1},
	}

	// SetUpsert(true) ensures that if no document matches the filter, a new one is created.
	opts := options.Update().SetUpsert(true)

	// Execute the atomic operation.
	_, err := s.indexCollection.UpdateOne(ctx, filter, update, opts)
	return err
}

const statsDocumentID = "global_stats"

// UpdateIndexStats updates the global statistics document.
// It uses an upsert operation to create the document if it does not exist.
func (s *MongoStore) UpdateIndexStats(ctx context.Context, totalDocs int64) error {
	// The filter targets the unique statistics document.
	filter := bson.M{"_id": statsDocumentID}

	// The update operation sets the total number of documents and the last update time.
	update := bson.M{
		"$set": bson.M{
			"total_documents": totalDocs,
			"last_indexed_at": time.Now(),
		},
	}

	// SetUpsert(true) ensures the document is created on the first run.
	opts := options.Update().SetUpsert(true)

	_, err := s.statsCollection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (s *MongoStore) GetIndexStats(ctx context.Context) (*IndexStats, error) {
	filter := bson.M{"_id": statsDocumentID}
	var stats IndexStats
	err := s.statsCollection.FindOne(ctx, filter).Decode(&stats)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &IndexStats{}, nil
		}
		return nil, err
	}
	return &stats, nil
}

// GetPostingsForTerms retrieves the inverted index entries for a given list of terms.
func (s *MongoStore) GetPostingsForTerms(ctx context.Context, terms []string) (map[string]InvertedIndexEntry, error) {
	if len(terms) == 0 {
		return make(map[string]InvertedIndexEntry), nil
	}

	filter := bson.M{"_id": bson.M{"$in": terms}}
	cursor, err := s.indexCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make(map[string]InvertedIndexEntry)
	for cursor.Next(ctx) {
		var entry InvertedIndexEntry
		if err := cursor.Decode(&entry); err != nil {
			return nil, err
		}
		results[entry.Term] = entry
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *MongoStore) GetDocuments(ctx context.Context, docIDs []string) ([]*Document, error) {
	if len(docIDs) == 0 {
		return []*Document{}, nil
	}

	objectIDs := make([]primitive.ObjectID, len(docIDs))
	for i, id := range docIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err // Invalid ID format
		}
		objectIDs[i] = objID
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	cursor, err := s.documentCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []*Document
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	return documents, nil
}

// BulkWriteDocuments performs a bulk write operation on the documents collection.
func (s *MongoStore) BulkWriteDocuments(ctx context.Context, models []mongo.WriteModel) error {
	if len(models) == 0 {
		return nil
	}
	opts := options.BulkWrite().SetOrdered(false) // Unordered for better performance
	_, err := s.documentCollection.BulkWrite(ctx, models, opts)
	return err
}

// BulkWriteInvertedIndex performs a bulk write operation on the inverted index collection.
func (s *MongoStore) BulkWriteInvertedIndex(ctx context.Context, models []mongo.WriteModel) error {
	if len(models) == 0 {
		return nil
	}
	opts := options.BulkWrite().SetOrdered(false)
	_, err := s.indexCollection.BulkWrite(ctx, models, opts)
	return err
}

// GetDocumentByPath retrieves a document by its file path.
func (s *MongoStore) GetDocumentByPath(ctx context.Context, filePath string) (*Document, error) {
	filter := bson.M{"file_path": filePath}
	var doc Document
	err := s.documentCollection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Return nil, nil when no document is found
		}
		return nil, err
	}
	return &doc, nil
}

// DeleteDocument deletes a document from the 'documents' collection by its ID.
func (s *MongoStore) DeleteDocument(ctx context.Context, docID primitive.ObjectID) error {
	_, err := s.documentCollection.DeleteOne(ctx, bson.M{"_id": docID})
	return err
}

// RemovePostingsForDocument removes all postings for a given document ID from the inverted index for a given list of terms.
func (s *MongoStore) RemovePostingsForDocument(ctx context.Context, docID primitive.ObjectID, terms []string) error {
	if len(terms) == 0 {
		return nil
	}

	filter := bson.M{"_id": bson.M{"_in": terms}}
	update := bson.M{"_pull": bson.M{"postings": bson.M{"doc_id": docID}}}

	_, err := s.indexCollection.UpdateMany(ctx, filter, update)
	return err
}
