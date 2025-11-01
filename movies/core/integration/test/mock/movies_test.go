package mock

import (
	"movies/core/proto"
	"movies/infra/persistence/mock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMoviesRepositoryMock_FindAll(t *testing.T) {
	mockRepo := mock.NewMoviesRepositoryMock()
	mock := mockRepo.(*mock.MoviesRepositoryMock)

	movies := []*proto.Movie{
		{Id: 1, Title: "The Matrix", Year: "1999"},
		{Id: 2, Title: "Inception", Year: "2010"},
		{Id: 3, Title: "Interstellar", Year: "2014"},
	}
	mock.Seed(movies)

	t.Run("should return paginated results", func(t *testing.T) {
		req := &proto.GetMoviesRequest{
			Page:  1,
			Limit: 2,
		}

		result, total, err := mockRepo.FindAll(req)

		require.NoError(t, err)
		assert.Equal(t, uint32(3), total)
		assert.Len(t, result, 2)

		assert.Equal(t, "Interstellar", result[0].Title)
		assert.Equal(t, "Inception", result[1].Title)
	})

	t.Run("should return empty for out of bounds page", func(t *testing.T) {
		req := &proto.GetMoviesRequest{
			Page:  10,
			Limit: 2,
		}

		result, total, err := mockRepo.FindAll(req)

		require.NoError(t, err)
		assert.Equal(t, uint32(3), total)
		assert.Len(t, result, 0)
	})
}

func TestMoviesRepositoryMock_FindById(t *testing.T) {
	mockRepo := mock.NewMoviesRepositoryMock()
	mock := mockRepo.(*mock.MoviesRepositoryMock)

	mock.Seed([]*proto.Movie{
		{Id: 1, Title: "The Matrix", Year: "1999"},
	})

	t.Run("should return movie when found", func(t *testing.T) {
		movie, err := mockRepo.FindById(&proto.MovieIdRequest{Id: 1})

		require.NoError(t, err)
		assert.Equal(t, "The Matrix", movie.Title)
		assert.Equal(t, "1999", movie.Year)
	})

	t.Run("should return error when not found", func(t *testing.T) {
		movie, err := mockRepo.FindById(&proto.MovieIdRequest{Id: 999})

		assert.Error(t, err)
		assert.Nil(t, movie)
	})
}

func TestMoviesRepositoryMock_Create(t *testing.T) {
	mockRepo := mock.NewMoviesRepositoryMock()

	t.Run("should create movie with auto-increment ID", func(t *testing.T) {
		newMovie := &proto.Movie{
			Title: "The Dark Knight",
			Year:  "2008",
		}

		created, err := mockRepo.Create(newMovie)

		require.NoError(t, err)
		assert.Equal(t, uint32(1), created.Id)
		assert.Equal(t, "The Dark Knight", created.Title)
		assert.Equal(t, "2008", created.Year)

		found, err := mockRepo.FindById(&proto.MovieIdRequest{Id: 1})
		require.NoError(t, err)
		assert.Equal(t, created, found)
	})

	t.Run("should increment ID for multiple creations", func(t *testing.T) {
		movie2 := &proto.Movie{Title: "Movie 2", Year: "2020"}
		movie3 := &proto.Movie{Title: "Movie 3", Year: "2021"}

		created2, err := mockRepo.Create(movie2)
		require.NoError(t, err)
		assert.Equal(t, uint32(2), created2.Id)

		created3, err := mockRepo.Create(movie3)
		require.NoError(t, err)
		assert.Equal(t, uint32(3), created3.Id)
	})
}

func TestMoviesRepositoryMock_Delete(t *testing.T) {
	mockRepo := mock.NewMoviesRepositoryMock()
	mock := mockRepo.(*mock.MoviesRepositoryMock)

	mock.Seed([]*proto.Movie{
		{Id: 1, Title: "The Matrix", Year: "1999"},
		{Id: 2, Title: "Inception", Year: "2010"},
	})

	t.Run("should delete existing movie", func(t *testing.T) {
		err := mockRepo.Delete(1)
		require.NoError(t, err)

		movie, err := mockRepo.FindById(&proto.MovieIdRequest{Id: 1})
		assert.Error(t, err)
		assert.Nil(t, movie)

		otherMovie, err := mockRepo.FindById(&proto.MovieIdRequest{Id: 2})
		require.NoError(t, err)
		assert.Equal(t, "Inception", otherMovie.Title)
	})

	t.Run("should return error for non-existent movie", func(t *testing.T) {
		err := mockRepo.Delete(999)
		assert.Error(t, err)
	})
}

func TestMoviesRepositoryMock_ClearAndSeed(t *testing.T) {
	mockRepo := mock.NewMoviesRepositoryMock()
	mock := mockRepo.(*mock.MoviesRepositoryMock)

	mockRepo.Create(&proto.Movie{Title: "Movie 1", Year: "2000"})
	mockRepo.Create(&proto.Movie{Title: "Movie 2", Year: "2001"})

	assert.Equal(t, uint32(3), mock.GetNextID())

	mock.Clear()

	movies, total, err := mockRepo.FindAll(&proto.GetMoviesRequest{Page: 1, Limit: 10})
	require.NoError(t, err)
	assert.Equal(t, uint32(0), total)
	assert.Len(t, movies, 0)
	assert.Equal(t, uint32(1), mock.GetNextID())

	seedMovies := []*proto.Movie{
		{Id: 10, Title: "Seeded 1", Year: "1990"},
		{Id: 20, Title: "Seeded 2", Year: "1995"},
	}
	mock.Seed(seedMovies)

	movie, err := mockRepo.FindById(&proto.MovieIdRequest{Id: 10})
	require.NoError(t, err)
	assert.Equal(t, "Seeded 1", movie.Title)
	assert.Equal(t, uint32(21), mock.GetNextID())
}
