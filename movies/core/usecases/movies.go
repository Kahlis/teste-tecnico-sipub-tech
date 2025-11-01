package usecases

import (
	"context"
	"movies/core/proto"
	"movies/core/repository"
	"movies/core/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MoviesUsecase struct {
	proto.UnimplementedMovieServiceServer
	Repository repository.MoviesRepository
}

func (service *MoviesUsecase) GetMovie(ctx context.Context, req *proto.MovieIdRequest) (*proto.Movie, error) {
	movie, err := service.Repository.FindById(req)

	if err == util.ErrMovieNotFound {
		return nil, status.Errorf(codes.NotFound, "movie not found")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch movie")
	}

	return movie, nil
}

func (service *MoviesUsecase) GetMovies(ctx context.Context, req *proto.GetMoviesRequest) (*proto.MovieListResponse, error) {
	movies, total, err := service.Repository.FindAll(req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch movies")
	}

	return &proto.MovieListResponse{
		Movies: movies,
		Total:  total,
		More:   false,
		Page:   req.Page,
	}, nil
}

func (service *MoviesUsecase) CreateMovie(ctx context.Context, req *proto.Movie) (*proto.Movie, error) {
	movie, err := service.Repository.Create(req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create movie")
	}

	return movie, nil
}

func (service *MoviesUsecase) DeleteMovie(ctx context.Context, req *proto.MovieIdRequest) (*proto.Empty, error) {
	err := service.Repository.Delete(req.Id)
	empty := &proto.Empty{}

	if err == util.ErrMovieNotFound {
		return empty, status.Errorf(codes.NotFound, "movie not found")
	}

	if err != nil {
		return empty, status.Errorf(codes.Internal, "failed to delete movie")
	}

	return empty, nil
}
