package util

import "errors"

var ErrMovieNotFound = errors.New("movie not found")
var ErrMovieAlreadyExists = errors.New("movie already exists")
