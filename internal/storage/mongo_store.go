package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	client            *mongo.Client
	database          *mongo.Database
	documetCollection *mongo.Collection
	indexCollection   *mongo.Collection
	statsCollection   *mongo.Collection
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
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")
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
	s.documetCollection = s.database.Collection("documents")
	s.indexCollection = s.database.Collection("inverted_index")
	s.statsCollection = s.database.Collection("stats")

	fmt.Println("Connected to MongoDB successfully.")
	return nil
}

// Disconnect safely closes the database connection.
func (s *MongoStore) Disconnect(ctx context.Context) error {
	// Logic to call s.client.Disconnect(ctx)
	// and handle the possible error.
	if err := s.client.Disconnect(ctx); err != nil {
		return err
	}
	fmt.Println("Disconnected from MongoDB.")
	// The actual implementation would go here...
	return nil // Placeholder
}

// AddDocument inserts a new document into the 'documents' collection.
// Returns the ID of the inserted document as a hexadecimal string.
func (s *MongoStore) AddDocument(ctx context.Context, doc Document) (string, error) {
	// Logic for:
	// 1. Call s.documentsCollection.InsertOne(ctx, doc)
	result, err := s.documetCollection.InsertOne(ctx, doc)
	// 2. Check for errors.
	if err != nil {
		return "", err
	}
	// 3. If there is no error, get the result.InsertedID.
	insertedID := result.InsertedID
	// 4. Type assert to primitive.ObjectID.
	objectID, ok := insertedID.(primitive.ObjectID)
	if !ok {
		return "", mongo.ErrInvalidIndexValue
	}
	// 5. Return objectID.Hex() and nil.
	return objectID.Hex(), nil
}

// UpsertTerm updates or inserts an entry in the inverted index.
// It searches for a term and adds the posting to its list. If the term does not exist,
// it creates it along with the initial posting. This operation is atomic.
func (s *MongoStore) UpsertTerm(ctx context.Context, term string, posting Posting) error {
	// Logic for:
	// 1. Define the `filter` (bson.M{"term": term}).
	filter := bson.M{"term": term}
	// 2. Define the `update` (bson.M{"$push": ..., "$setOnInsert": ...}).
	update := bson.M{
		"$push": bson.M{
			"postings": posting,
		},
		"$setOnInsert": bson.M{
			"term": term,
		},
	}
	// 3. Define the `options` (options.Update().SetUpsert(true)).
	updateOptions := options.Update().SetUpsert(true)
	// 4. Call s.indexCollection.UpdateOne(ctx, filter, update, opts).
	_, err := s.indexCollection.UpdateOne(ctx, filter, update, updateOptions)
	if err != nil {
		return err
	}

	return nil
}

// UpdateIndexStats updates the global statistics document.
// It uses an upsert operation to create the document if it does not exist.
func (s *MongoStore) UpdateIndexStats(ctx context.Context, totalDocs int64) error {
	// Logic for:
	// 1. Define a constant `filter` for the stats document.
	filter := bson.M{} // Assuming there is only one stats document
	// 2. Define the `update` (bson.M{"$set": ...}).
	update := bson.M{
		"$set": bson.M{
			"total_documents": totalDocs,
			"last_indexed_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	// 3. Define the `options` (options.Update().SetUpsert(true)).
	opts := options.Update().SetUpsert(true)
	// 4. Call s.statsCollection.UpdateOne(ctx, filter, update, opts).
	_, err := s.statsCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	// 5. Return the resulting error.

	// The actual implementation would go here...
	return nil // Placeholder
}
