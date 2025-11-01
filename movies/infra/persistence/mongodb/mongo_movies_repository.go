package mongodb

import (
	"context"
	"movies/core/proto"
	"movies/core/repository"
	"movies/core/util"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MoviesRepositoryImpl struct {
	collection *mongo.Collection
}

func NewMoviesRepository(client *mongo.Client, dbName, collectionName string) repository.MoviesRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &MoviesRepositoryImpl{collection: collection}
}

func (repo *MoviesRepositoryImpl) FindAll(req *proto.GetMoviesRequest) ([]*proto.Movie, uint32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().
		SetSkip(int64(req.Page-1) * int64(req.Limit)).
		SetLimit(int64(req.Limit)).
		SetSort(bson.M{"createdAt": -1})

	cursor, err := repo.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var movies []*proto.Movie
	if err = cursor.All(ctx, &movies); err != nil {
		return nil, 0, err
	}

	total, err := repo.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return movies, uint32(total), nil
}

func (repo *MoviesRepositoryImpl) FindById(req *proto.MovieIdRequest) (*proto.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var movie proto.Movie
	repo.collection.FindOne(ctx, bson.M{"id": req.Id}).Decode(&movie)
	if movie.Id != req.Id {
		return nil, util.ErrMovieNotFound
	}

	return &movie, nil
}

func (repo *MoviesRepositoryImpl) Create(movie *proto.Movie) (*proto.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	movie.Id = GetNextID()
	_, err := repo.collection.InsertOne(ctx, movie)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (repo *MoviesRepositoryImpl) Delete(id uint32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := repo.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return util.ErrMovieNotFound
	}

	return nil
}
