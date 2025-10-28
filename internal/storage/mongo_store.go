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

// NewMongoStore es una función constructora que inicializa una instancia de MongoStore.
// La conexión real se establece a través del método Connect().
func NewMongoStore() *MongoStore {
	return &MongoStore{}
}

// Connect establece la conexión con MongoDB usando una URI de una variable de entorno.
// También realiza un ping para verificar la conexión y prepara los manejadores de las colecciones.
func (s *MongoStore) Connect(ctx context.Context) error {
	// Lógica detallada para:
	// 1. Leer os.Getenv("MONGO_URI")
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")
	// 2. Configurar clientOptions
	clientOptions := options.Client().ApplyURI(mongoURI)
	// 3. Llamar a mongo.Connect(ctx, clientOptions)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	// 4. Llamar a client.Ping(ctx, nil) para verificar
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}
	// 5. Si todo va bien, inicializar los campos de la struct:
	s.client = client
	s.database = client.Database(dbName)
	s.documetCollection = s.database.Collection("documents")
	s.indexCollection = s.database.Collection("inverted_index")
	s.statsCollection = s.database.Collection("stats")

	fmt.Println("Connected to MongoDB successfully.")
	return nil
}

// Disconnect cierra la conexión con la base de datos de forma segura.
func (s *MongoStore) Disconnect(ctx context.Context) error {
	// Lógica para llamar a s.client.Disconnect(ctx)
	// y manejar el posible error.
	if err := s.client.Disconnect(ctx); err != nil {
		return err
	}
	fmt.Println("Disconnected from MongoDB.")
	// La implementación real iría aquí...
	return nil // Placeholder
}

// AddDocument inserta un nuevo documento en la colección 'documents'.
// Devuelve el ID del documento insertado como una string hexadecimal.
func (s *MongoStore) AddDocument(ctx context.Context, doc Document) (string, error) {
	// Lógica para:
	// 1. Llamar a s.documentsCollection.InsertOne(ctx, doc)
	result, err := s.documetCollection.InsertOne(ctx, doc)
	// 2. Verificar el error.
	if err != nil {
		return "", err
	}
	// 3. Si no hay error, obtener el result.InsertedID.
	insertedID := result.InsertedID
	// 4. Hacer una aserción de tipo a primitive.ObjectID.
	objectID, ok := insertedID.(primitive.ObjectID)
	if !ok {
		return "", mongo.ErrInvalidIndexValue
	}
	// 5. Devolver objectID.Hex() y nil.
	return objectID.Hex(), nil
}

// UpsertTerm actualiza o inserta una entrada en el índice invertido.
// Busca un término y añade el posting a su lista. Si el término no existe,
// lo crea junto con el posting inicial. Esta operación es atómica.
func (s *MongoStore) UpsertTerm(ctx context.Context, term string, posting Posting) error {
	// Lógica para:
	// 1. Definir el `filter` (bson.M{"term": term}).
	filter := bson.M{"term": term}
	// 2. Definir el `update` (bson.M{"$push": ..., "$setOnInsert": ...}).
	update := bson.M{
		"$push": bson.M{
			"postings": posting,
		},
		"$setOnInsert": bson.M{
			"term": term,
		},
	}
	// 3. Definir las `options` (options.Update().SetUpsert(true)).
	updateOptions := options.Update().SetUpsert(true)
	// 4. Llamar a s.indexCollection.UpdateOne(ctx, filter, update, opts).
	_, err := s.indexCollection.UpdateOne(ctx, filter, update, updateOptions)
	if err != nil {
		return err
	}

	return nil
}

// UpdateIndexStats actualiza el documento de estadísticas globales.
// Utiliza una operación de upsert para crear el documento si no existe.
func (s *MongoStore) UpdateIndexStats(ctx context.Context, totalDocs int64) error {
	// Lógica para:
	// 1. Definir un `filter` constante para el documento de stats.
	filter := bson.M{} // Asumiendo que solo hay un documento de stats
	// 2. Definir el `update` (bson.M{"$set": ...}).
	update := bson.M{
		"$set": bson.M{
			"total_documents": totalDocs,
			"last_indexed_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	// 3. Definir las `options` (options.Update().SetUpsert(true)).
	opts := options.Update().SetUpsert(true)
	// 4. Llamar a s.statsCollection.UpdateOne(ctx, filter, update, opts).
	_, err := s.statsCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	// 5. Devolver el error resultante.

	// La implementación real iría aquí...
	return nil // Placeholder
}
