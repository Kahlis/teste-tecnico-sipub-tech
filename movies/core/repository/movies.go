package repository

import (
	"movies/core/proto"
)

type MoviesRepositoryImpl struct {
	movies map[uint32]*proto.Movie
	nextID uint32
}

type MoviesRepository interface {
	FindAll(req *proto.GetMoviesRequest) ([]*proto.Movie, uint32, error)
	FindById(req *proto.MovieIdRequest) (*proto.Movie, error)
	Create(movie *proto.Movie) (*proto.Movie, error)
	Delete(id uint32) error
}
