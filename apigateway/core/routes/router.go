package routes

import (
	"apigateway/core/handler"
	"apigateway/core/usecases"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Register(
	router *gin.Engine,
	moviesUsecase *usecases.MoviesUsecases,
	logger *zap.Logger,
) {
	api := router.Group("/v1")
	{
		handler.RegisterMoviesRoutes(api, moviesUsecase, logger)
		handler.RegisterHealthRoute(api)
		handler.RegisterSwagger(api)
	}
}
