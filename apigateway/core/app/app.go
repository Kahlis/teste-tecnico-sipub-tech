package app

import (
	"apigateway/core/config"
	"apigateway/core/routes"
	"apigateway/core/usecases"
	"apigateway/infra/clients"
	"apigateway/pkg/logger"
	"fmt"

	_ "apigateway/docs"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title           API Gateway Swagger
// @version         1.0
// @description     API Gateway for Movies Service
// @host            localhost:8080
// @BasePath        /
// @schemes         http
func SetupRouter() (*gin.Engine, *zap.Logger, error) {
	cfg := config.Load()
	log := logger.New(cfg.Env)

	grpcClient, err := clients.NewGrpcClient(&cfg)

	if err != nil {
		return nil, log, err
	}

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	moviesUsecases := usecases.NewMoviesUseCases(grpcClient, log)
	router := gin.Default()
	router.SetTrustedProxies(nil)
	routes.Register(router, moviesUsecases, log)

	return router, log, nil
}

func Run() error {
	router, log, err := SetupRouter()
	if err != nil {
		return err
	}

	cfg := config.Load()
	listenPort := fmt.Sprintf(":%s", cfg.ListenPort)

	log.Info("API Gateway running", zap.String("address", listenPort))
	return router.Run(listenPort)
}
