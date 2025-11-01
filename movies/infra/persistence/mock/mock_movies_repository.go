package mock

import (
	"movies/core/proto"
	"movies/core/repository"
	"movies/core/util"
	"sort"
)

type MoviesRepositoryMock struct {
	movies map[uint32]*proto.Movie
	nextID uint32
}

type Mock interface {
	Clear()
	Seed([]*proto.Movie)
	GetNextID() uint32
}

var _ Mock = (*MoviesRepositoryMock)(nil)

func NewMoviesRepositoryMock() repository.MoviesRepository {
	return &MoviesRepositoryMock{
		movies: make(map[uint32]*proto.Movie),
		nextID: 1,
	}
}

func (repo *MoviesRepositoryMock) FindAll(req *proto.GetMoviesRequest) ([]*proto.Movie, uint32, error) {
	movies := make([]*proto.Movie, 0, len(repo.movies))
	for _, movie := range repo.movies {
		movies = append(movies, movie)
	}

	sort.Slice(movies, func(i, j int) bool {
		return movies[i].Id > movies[j].Id
	})

	start := (req.Page - 1) * req.Limit
	if start >= uint32(len(movies)) {
		return []*proto.Movie{}, uint32(len(movies)), nil
	}

	end := start + req.Limit
	if end > uint32(len(movies)) {
		end = uint32(len(movies))
	}

	return movies[start:end], uint32(len(movies)), nil
}

func (repo *MoviesRepositoryMock) FindById(req *proto.MovieIdRequest) (*proto.Movie, error) {
	movie, exists := repo.movies[req.Id]
	if !exists {
		return nil, util.ErrMovieNotFound
	}
	return movie, nil
}

func (repo *MoviesRepositoryMock) Create(movie *proto.Movie) (*proto.Movie, error) {
	movie.Id = repo.nextID
	repo.movies[movie.Id] = movie
	repo.nextID++
	return movie, nil
}

func (repo *MoviesRepositoryMock) Delete(id uint32) error {
	if _, exists := repo.movies[id]; !exists {
		return util.ErrMovieNotFound
	}
	delete(repo.movies, id)
	return nil
}
func (repo *MoviesRepositoryMock) Clear() {
	repo.movies = make(map[uint32]*proto.Movie)
	repo.nextID = 1
}

func (repo *MoviesRepositoryMock) Seed(movies []*proto.Movie) {
	for _, movie := range movies {
		repo.movies[movie.Id] = movie
		if movie.Id >= repo.nextID {
			repo.nextID = movie.Id + 1
		}
	}
}

func (repo *MoviesRepositoryMock) GetNextID() uint32 {
	return repo.nextID
}
