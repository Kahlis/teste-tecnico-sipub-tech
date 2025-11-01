package util

import (
	"errors"
)

const (
	pageNumberInvalid = "page number must be greater than 0"
	pageSizeShort     = "page size must be greater than 1"
	pageSizeLong      = "page size must be less than 20"
	movieNotFound     = "movie not found"
	moviePageNotFound = "movie page not found"
	titleEmpty        = "title cannot be empty"
	yearEmpty         = "year cannot be empty"
)

var ErrPageNumberInvalid = errors.New(pageNumberInvalid)
var ErrPageSizeShort = errors.New(pageSizeShort)
var ErrPageSizeLong = errors.New(pageSizeLong)
var ErrMovieNotFound = errors.New(movieNotFound)
var ErrMoviePageNotFound = errors.New(moviePageNotFound)
var ErrTitleEmpty = errors.New(titleEmpty)
var ErrYearEmpty = errors.New(yearEmpty)

func IsErrInvalidParams(err error) bool {
	if err == ErrPageNumberInvalid || err == ErrPageSizeShort || err == ErrPageSizeLong {
		return true
	}
	return false
}

func IsInvalidBody(err error) bool {
	if err == ErrTitleEmpty || err == ErrYearEmpty {
		return true
	}
	return false
}
