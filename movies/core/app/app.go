package app

import (
	"fmt"
	"movies/core/config"
	"movies/core/proto"
	"movies/core/usecases"
	"movies/infra/persistence/mongodb"
	"movies/pkg/logger"
	"net"

	"google.golang.org/grpc"
)

func Run() error {
	cfg := config.Load()
	log := logger.New(cfg.Env)

	grpcServer := grpc.NewServer()

	mongodb.SetLogger(log)
	db, err := mongodb.ConnectToMongo(&cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Mongo connection was setup")
	movies := mongodb.NewMoviesRepository(db, cfg.DbName, cfg.DbCollection)
	proto.RegisterMovieServiceServer(grpcServer, &usecases.MoviesUsecase{Repository: movies})

	listenPort := fmt.Sprintf(":%s", cfg.ListenPort)
	listener, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatal(err.Error())
	}

	logInfo := fmt.Sprintf("Movies services running on port %s", listenPort)
	log.Info(logInfo)

	grpcErr := grpcServer.Serve(listener)
	if grpcErr != nil {
		log.Fatal(grpcErr.Error())
	}

	return nil
}
