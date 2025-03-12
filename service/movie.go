package service

import (
	"github.com/rodrigo462003/FlickMeter/model"
	"github.com/rodrigo462003/FlickMeter/movieAPI"
)

type MovieService interface {
	GetMovie(movieID string) (movie *model.Movie, err error)
	SearchMovies(search string) (movies []model.MovieIndex)
}

type movieService struct {
	movieAPI.MovieGetter
}

func NewMovieService(apiToken, filePath string) *movieService {
	return &movieService{movieAPI.NewMovieGet(apiToken, filePath)}
}

func (s *movieService) GetMovie(movieID string) (movie *model.Movie, err error) {
	movie, err = s.MovieGetter.GetMovie(movieID)
	if err != nil {
		return nil, err
	}

	videos, err := s.MovieGetter.GetVideos(movieID)
	if err != nil {
		return nil, err
	}
	movie.Videos = videos

	return movie, nil
}

func (s *movieService) SearchMovies(search string) (movies []model.MovieIndex) {
	movies = s.MovieGetter.MoviesIndex()
	return movies[len(movies)-5:]
}
