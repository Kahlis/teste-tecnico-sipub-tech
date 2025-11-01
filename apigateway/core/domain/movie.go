package domain

import (
	"apigateway/core/proto"
	"apigateway/core/util"
)

type Movie struct {
	Id    uint32 `json:"id"`
	Title string `json:"title"`
	Year  string `json:"year"`
}

type MovieList struct {
	Movies  []*Movie `json:"movies"`
	More    bool     `json:"more"`
	Page    uint32   `json:"page"`
	Total   uint32   `json:"total"`
	Results uint32   `json:"results"`
}

func ParseMovie(movie *proto.Movie) *Movie {
	return &Movie{
		Id:    movie.Id,
		Title: movie.Title,
		Year:  movie.Year,
	}
}

func ParseMovies(movies []*proto.Movie) []*Movie {
	var parsedMovies []*Movie
	for _, movie := range movies {
		parsedMovie := ParseMovie(movie)
		parsedMovies = append(parsedMovies, parsedMovie)
	}

	return parsedMovies
}

func ParseMovieList(movies *proto.MovieListResponse) *MovieList {
	parsedMovies := ParseMovies(movies.Movies)

	return &MovieList{
		Movies: parsedMovies,
		More:   movies.More,
		Page:   movies.Page,
		Total:  movies.Total,
	}
}

func HasMore(total uint32, page uint32, resultsPerPage uint32) bool {
	return total > page*resultsPerPage
}

func IsPageNumberValid(pageNumber int) error {
	if pageNumber < 1 {
		return util.ErrPageNumberInvalid
	}

	return nil
}

func IsResultsPerPageValid(pageSize int) error {
	if pageSize < 2 {
		return util.ErrPageSizeShort
	}

	if pageSize > 20 {
		return util.ErrPageSizeLong
	}

	return nil
}

func IsValidMovie(movie *Movie) error {
	if movie.Title == "" {
		return util.ErrTitleEmpty
	}

	if movie.Year == "" {
		return util.ErrYearEmpty
	}

	return nil
}
