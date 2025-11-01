package mongodb

import (
	"context"
	"encoding/json"
	"movies/core/config"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var collection *mongo.Collection
var log *zap.Logger
var currentID int = 0

func SetLogger(logger *zap.Logger) {
	log = logger
}

func ConnectToMongo(cfg *config.Config) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(cfg.DBUri)

	clientOptions.SetAuth(options.Credential{
		Username: cfg.DBUser,
		Password: cfg.DBPassword,
	})

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	collection = client.Database(cfg.DbName).Collection(cfg.DbCollection)

	uniqueErr := createUniqueIndexes()
	if uniqueErr != nil {
		log.Error("Failed to create indexes", zap.Error(uniqueErr))
		return nil, uniqueErr
	}

	count, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}

	if count == 0 {
		log.Info("Seeding from JSON file...")
		if err := SeedFromJSON(context.Background(), collection); err != nil {
			return nil, err
		}
	}

	initializeIdErr := initializeCurrentID()
	if initializeIdErr != nil {
		log.Error("Failed to initialize current ID", zap.Error(initializeIdErr))
		return nil, initializeIdErr
	}

	return client, nil
}

func GetCollection() *mongo.Collection {
	return collection
}

func createUniqueIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("unique_id"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

func initializeCurrentID() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "maxId", Value: bson.D{{Key: "$max", Value: "$id"}}},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result struct {
			MaxID int `bson:"maxId"`
		}
		if err := cursor.Decode(&result); err != nil {
			return err
		}
		currentID = result.MaxID
	} else {
		currentID = 0
	}

	return nil
}

func GetNextID() uint32 {
	currentID++
	return uint32(currentID)
}

func SeedFromJSON(ctx context.Context, collection *mongo.Collection) error {
	data, err := os.ReadFile(GetSeedPath())
	if err != nil {
		return err
	}

	var movies []interface{}
	if err := json.Unmarshal(data, &movies); err != nil {
		return err
	}

	if len(movies) > 0 {
		_, err = collection.InsertMany(ctx, movies)
		return err
	}

	return nil
}

func GetSeedPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "..", "..", "seed", "movies.json")
}
