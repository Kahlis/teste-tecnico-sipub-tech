package e2e

import (
	"encoding/json"
	"fmt"
	"movies/core/config"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseUrl string

type Post struct {
	Title string `json:"title"`
	Year  string `json:"year"`
}

type PostBadRequest struct {
	Title string `json:"title"`
}

func requiredField(field string, found ...string) string {
	if len(found) > 0 {
		return fmt.Sprintf("response should contain `%s` field, found `%s` instead", field, found[0])
	}
	return fmt.Sprintf("response should contain `%s` field", field)
}

func shouldNotBeError(test *testing.T, res, route string) {
	if strings.Contains(res, (`"error"`)) {
		fatalMessage := fmt.Sprintf(
			"Response from `%s` should not be an error",
			route,
		)
		test.Fatal(fatalMessage)
	}
}

func TestMain(m *testing.M) {
	cfg := config.Load()
	baseUrl = fmt.Sprintf("http://localhost:%s", cfg.ApiPort)

	os.Exit(m.Run())
}

func TestCreateMovie(test *testing.T) {
	route := "/v1/movies"
	tc := NewTestClient(test, baseUrl)

	payload, err := json.Marshal(&Post{"Alice in Wonderland", "2010"})

	if err != nil {
		test.Fatal(err)
	}

	res := tc.Post(route, payload)
	assert.Equal(test, http.StatusCreated, res.StatusCode)

	// Ensure response default pattern
	assert.Contains(test, res.Body, `"data"`, "response should contain `data` field")
	assert.Contains(test, res.Body, `"success"`, "response should contain `success` field")
	assert.Contains(test, res.Body, `"success":true`, "`success` field should be false")

	// Ensure specific fields for Movie
	assert.Contains(test, res.Body, `"id"`, "response should contain `id` field")
	assert.Contains(test, res.Body, `"title"`, "response should contain `title` field")
	assert.Contains(test, res.Body, `"year"`, "response should contain `year` field")
}

func TestGetMovie(test *testing.T) {
	route := "/v1/movies/7087851"
	tc := NewTestClient(test, baseUrl)

	res := tc.Get(route)
	shouldNotBeError(test, res.Body, route)
	assert.Equal(test, http.StatusOK, res.StatusCode)

	// Ensure response default pattern
	assert.Contains(test, res.Body, `"data"`, "response should contain `data` field")
	assert.Contains(test, res.Body, `"success"`, "response should contain `success` field")
	assert.Contains(test, res.Body, `"success":true`, "`success` field should be false")

	// Ensure specific fields for Movie
	assert.Contains(test, res.Body, `"id"`, "response should contain `id` field")
	assert.Contains(test, res.Body, `"title"`, requiredField("title"))
	assert.Contains(test, res.Body, `"year"`, requiredField("year"))
}

func TestGetMovieList(test *testing.T) {
	route := "/v1/movies"
	tc := NewTestClient(test, baseUrl)

	res := tc.Get(route)
	assert.Equal(test, http.StatusOK, res.StatusCode)

	keyRegex := regexp.MustCompile(`"([^"]+)":\s*\[.*?\]`)
	keyMatch := keyRegex.FindStringSubmatch(res.Body)

	shouldNotBeError(test, res.Body, route)

	// Ensure MovieList Movies
	if len(keyMatch) < 2 {
		fatalMessage := fmt.Sprintf(
			"Failed to find `movies` field in response from `%s`",
			route,
		)
		test.Fatal(fatalMessage)
	} else {
		assert.Regexp(
			test,
			`"movies":\[.*?\]`,
			res.Body,
			requiredField("movies", keyMatch[1]),
		)
	}

	assert.Equal(test, http.StatusOK, res.StatusCode)

	// Ensure response default pattern
	assert.Contains(test, res.Body, `"data"`, "response should contain `data` field")
	assert.Contains(test, res.Body, `"success"`, "response should contain `success` field")
	assert.Contains(test, res.Body, `"success":true`, "`success` field should be false")

	// Ensure specific fields for MovieList
	assert.Contains(test, res.Body, `"more"`, requiredField("more"))
	assert.Contains(test, res.Body, `"page"`, requiredField("page"))
	assert.Contains(test, res.Body, `"total"`, requiredField("total"))
	assert.Contains(test, res.Body, `"results"`, requiredField("results"))
}

func TestDeleteMovie(test *testing.T) {
	route := "/v1/movies/7087851"
	tc := NewTestClient(test, baseUrl)
	res := tc.Delete(route)

	shouldNotBeError(test, res.Body, route)
	assert.Equal(test, http.StatusNoContent, res.StatusCode)
}

func TestDeleteMovieNotFound(test *testing.T) {
	route := "/v1/movies/9999999"
	tc := NewTestClient(test, baseUrl)
	res := tc.Delete(route)

	assert.Equal(test, http.StatusNotFound, res.StatusCode)

	// Ensure error response default pattern
	assert.Contains(test, res.Body, `"error"`, requiredField("error"))
	assert.Contains(test, res.Body, `"success"`, requiredField("success"))
	assert.Contains(test, res.Body, `"success":false`, "`success` field should be false")
}

func TestCreateMovieBadRequest(test *testing.T) {
	route := "/v1/movies"
	tc := NewTestClient(test, baseUrl)

	payload, err := json.Marshal(&PostBadRequest{"Alice in Wonderland"})

	if err != nil {
		test.Fatal(err)
	}

	res := tc.Post(route, payload)

	assert.Equal(test, http.StatusBadRequest, res.StatusCode)

	// Ensure error response default pattern
	assert.Contains(test, res.Body, `"error"`, requiredField("error"))
	assert.Contains(test, res.Body, `"success"`, requiredField("success"))
	assert.Contains(test, res.Body, `"success":false`, "`success` field should be false")
}

func TestGetMovieNotFound(test *testing.T) {
	route := "/v1/movies/9999999"
	tc := NewTestClient(test, baseUrl)

	res := tc.Get(route)
	assert.Equal(test, http.StatusNotFound, res.StatusCode)

	// Ensure error response default pattern
	assert.Contains(test, res.Body, `"error"`, requiredField("error"))
	assert.Contains(test, res.Body, `"success"`, requiredField("success"))
	assert.Contains(test, res.Body, `"success":false`, "`success` field should be false")
}
