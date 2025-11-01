package usecases

import (
	"apigateway/core/domain"
	"apigateway/core/proto"
	"apigateway/core/util"
	"apigateway/infra/clients"
	"context"

	"go.uber.org/zap"
)

type MoviesUsecases struct {
	Client *clients.MoviesGRPCClient
	Logger *zap.Logger
}

type MoviesClient interface {
	GetMovie(ctx context.Context, id int) (*domain.Movie, error)
	GetMovies(ctx context.Context, pageNumber, resultsPerPage int) ([]*domain.Movie, error)
	CreateMovie(ctx context.Context, movie *domain.Movie) (domain.Movie, error)
	DeleteMovie(ctx context.Context, id int) error
}

func NewMoviesUseCases(client *clients.MoviesGRPCClient, logger *zap.Logger) *MoviesUsecases {
	return &MoviesUsecases{
		Client: client,
		Logger: logger,
	}
}

func (m *MoviesUsecases) GetMovie(ctx context.Context, id int) (*domain.Movie, error) {
	movieQuery, err := m.Client.GetMovie(ctx, uint32(id))

	if err != nil {
		return nil, err
	}

	movie := domain.ParseMovie(movieQuery)
	return movie, err
}

func (m *MoviesUsecases) GetMovies(ctx context.Context, pageNumber, resultsPerPage int) (*domain.MovieList, error) {
	err := domain.IsPageNumberValid(pageNumber)

	if err != nil {
		return nil, err
	}

	err = domain.IsResultsPerPageValid(resultsPerPage)

	if err != nil {
		return nil, err
	}

	movieList, err := m.Client.GetMovies(ctx, pageNumber, resultsPerPage)

	if err != nil {
		return nil, err
	}

	movies := domain.ParseMovies(movieList.Movies)
	if movies == nil {
		return nil, util.ErrMoviePageNotFound
	}

	hasMore := domain.HasMore(movieList.Total, movieList.Page, uint32(resultsPerPage))
	return &domain.MovieList{
		Movies:  movies,
		More:    hasMore,
		Total:   movieList.Total,
		Page:    movieList.Page,
		Results: uint32(resultsPerPage),
	}, err
}

func (m *MoviesUsecases) CreateMovie(ctx context.Context, movie *proto.Movie) (*domain.Movie, error) {

	movieQuery, err := m.Client.CreateMovie(ctx, movie)

	if err != nil {
		return nil, err
	}

	parsedMovie := domain.ParseMovie(movieQuery)

	if err := domain.IsValidMovie(parsedMovie); err != nil {
		return nil, err
	}

	return parsedMovie, nil
}

func (m *MoviesUsecases) DeleteMovie(ctx context.Context, id int) error {
	return m.Client.DeleteMovie(ctx, id)
}
