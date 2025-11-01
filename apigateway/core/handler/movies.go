package handler

import (
	"apigateway/core/proto"
	"apigateway/core/usecases"
	"apigateway/core/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

const (
	movieNotFoundMessage         = "movie not found"
	internalServerErrorMessage   = "internal server error"
	invalidPageNumberMessage     = "Invalid `pageNumber` value"
	invalidResultsPerPageMessage = "Invalid `resultsPerPage` value"
	invalidIdMessage             = "Invalid `id` value"
	invalidRequestMessage        = "Invalid request"
	pageNotFoundMessage          = "Page not found"
)

type MoviesHandler struct {
	UseCases *usecases.MoviesUsecases
	Logger   *zap.Logger
}

// @Summary Listar filmes
// @Description Retorna uma lista paginada de filmes cadastrados
// @Tags Movies
// @Accept json
// @Produce json
// @Param pageNumber query int false "Número da página (padrão 1)"
// @Param resultsPerPage query int false "Resultados por página (padrão 10)"
// @Success 200 {object} map[string]interface{} "Lista de filmes retornada com sucesso"
// @Failure 400 {object} map[string]interface{} "Parâmetro inválido"
// @Failure 404 {object} map[string]interface{} "Página não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /movies [get]
func (handler *MoviesHandler) GetMovies(context *gin.Context) {
	pageNumber := context.DefaultQuery("pageNumber", "1")
	resultsPerPage := context.DefaultQuery("resultsPerPage", "10")

	pageNumberInt, errPageNumber := strconv.Atoi(pageNumber)
	if errPageNumber != nil {
		util.SendError(context, http.StatusBadRequest, invalidPageNumberMessage, errPageNumber)
		return
	}

	resultsPerPageInt, errResultsPerPage := strconv.Atoi(resultsPerPage)
	if errResultsPerPage != nil {
		util.SendError(context, http.StatusBadRequest, invalidResultsPerPageMessage, errResultsPerPage)
		return
	}

	movies, err := handler.UseCases.GetMovies(context, pageNumberInt, resultsPerPageInt)

	if err != nil && err == util.ErrMoviePageNotFound {
		util.SendError(context, http.StatusNotFound, pageNotFoundMessage, err)
		return
	}

	if err != nil && util.IsErrInvalidParams(err) {
		util.SendError(context, http.StatusBadRequest, invalidRequestMessage, err)
		return
	}

	if err != nil {
		handler.Logger.Error("could not get movies", zap.Error(err))
		util.SendError(context, http.StatusInternalServerError, internalServerErrorMessage, err)
		return
	}

	util.SendSuccess(context, http.StatusOK, movies)
}

// @Summary Obter filme por ID
// @Description Retorna um único filme com base em seu ID
// @Tags Movies
// @Accept json
// @Produce json
// @Param id path int true "ID do filme"
// @Success 200 {object} map[string]interface{} "Filme retornado com sucesso"
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 404 {object} map[string]interface{} "Filme não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /movies/{id} [get]
func (handler *MoviesHandler) GetMovie(context *gin.Context) {
	id := context.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		util.SendError(context, http.StatusBadRequest, invalidIdMessage, err)
		return
	}

	movie, err := handler.UseCases.GetMovie(context, idInt)

	grpcErr := util.ParseGRPCError(err)

	if grpcErr != nil && grpcErr.Code == codes.NotFound {
		util.SendError(context, http.StatusNotFound, movieNotFoundMessage, err)
		return
	}

	if err != nil {
		code, message, details := util.GRPCToZap(grpcErr)
		handler.Logger.Error("Internal Server Error", code, message, details)
		util.SendError(context, http.StatusInternalServerError, internalServerErrorMessage, err)
		return
	}

	util.SendSuccess(context, http.StatusOK, movie)
}

func (handler *MoviesHandler) CreateMovie(context *gin.Context) {
	var movie proto.Movie

	if err := context.ShouldBindJSON(&movie); err != nil {
		util.SendError(context, http.StatusBadRequest, invalidRequestMessage, err)
		return
	}

	created, err := handler.UseCases.CreateMovie(context, &movie)

	if err != nil && util.IsInvalidBody(err) {
		util.SendError(context, http.StatusBadRequest, internalServerErrorMessage, err)
		return
	}

	if err != nil {
		handler.Logger.Error("Internal Server Error", zap.Error(err))
		util.SendError(context, http.StatusInternalServerError, internalServerErrorMessage, err)
		return
	}

	util.SendSuccess(context, http.StatusCreated, created)
}

// @Summary Criar novo filme
// @Description Registra um novo filme no banco de dados
// @Tags Movies
// @Accept json
// @Produce json
// @Param movie body proto.Movie true "Objeto do filme a ser criado"
// @Success 201 {object} map[string]interface{} "Filme criado com sucesso"
// @Failure 400 {object} map[string]interface{} "Requisição inválida"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /movies [post]
func (handler *MoviesHandler) DeleteMovie(context *gin.Context) {
	id := context.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		util.SendError(context, http.StatusBadRequest, invalidIdMessage, err)
		return
	}

	errDelete := handler.UseCases.DeleteMovie(context, idInt)

	grpcErr := util.ParseGRPCError(errDelete)

	if grpcErr != nil && grpcErr.Code == codes.NotFound {
		util.SendError(context, http.StatusNotFound, movieNotFoundMessage, errDelete)
		return
	}

	if errDelete != nil {
		code, message, details := util.GRPCToZap(grpcErr)
		handler.Logger.Error("Internal Server Error", code, message, details)
		util.SendError(context, http.StatusInternalServerError, movieNotFoundMessage, errDelete)
		return
	}

	util.SendSuccess(context, http.StatusNoContent, nil)
}

// @Summary Deletar filme por ID
// @Description Remove um filme do banco de dados
// @Tags Movies
// @Accept json
// @Produce json
// @Param id path int true "ID do filme"
// @Success 204 {object} nil "Filme deletado com sucesso"
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 404 {object} map[string]interface{} "Filme não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /movies/{id} [delete]
func RegisterMoviesRoutes(rg *gin.RouterGroup, moviesUseCases *usecases.MoviesUsecases, logger *zap.Logger) {
	moviesHandler := &MoviesHandler{
		UseCases: moviesUseCases,
		Logger:   logger,
	}

	movies := rg.Group("/movies")

	movies.GET("", moviesHandler.GetMovies)          // List all movies
	movies.GET("/:id", moviesHandler.GetMovie)       // Get a movie by ID
	movies.POST("", moviesHandler.CreateMovie)       // Create a new movie
	movies.DELETE("/:id", moviesHandler.DeleteMovie) // Delete a movie by ID
}
