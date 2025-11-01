package clients

import (
	"apigateway/core/config"
	"apigateway/core/proto"
	"context"

	"google.golang.org/grpc"
)

type MoviesGRPCClient struct {
	Conn   *grpc.ClientConn
	Client proto.MovieServiceClient
}

func NewGrpcClient(cfg *config.Config) (*MoviesGRPCClient, error) {
	conn, err := grpc.Dial(cfg.GrpcServerAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := proto.NewMovieServiceClient(conn)

	return &MoviesGRPCClient{
		Conn:   conn,
		Client: client,
	}, nil
}

func (c *MoviesGRPCClient) Close() error {
	return c.Conn.Close()
}

func (c *MoviesGRPCClient) GetMovie(ctx context.Context, id uint32) (*proto.Movie, error) {
	resp, err := c.Client.GetMovie(ctx, &proto.MovieIdRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *MoviesGRPCClient) GetMovies(ctx context.Context, page, results int) (*proto.MovieListResponse, error) {
	resp, err := c.Client.GetMovies(ctx, &proto.GetMoviesRequest{Page: uint32(page), Limit: uint32(results)})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *MoviesGRPCClient) CreateMovie(ctx context.Context, movie *proto.Movie) (*proto.Movie, error) {
	resp, err := c.Client.CreateMovie(ctx, &proto.Movie{
		Title: movie.Title,
		Year:  movie.Year,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *MoviesGRPCClient) DeleteMovie(ctx context.Context, id int) error {
	_, err := c.Client.DeleteMovie(ctx, &proto.MovieIdRequest{Id: uint32(id)})
	return err
}
