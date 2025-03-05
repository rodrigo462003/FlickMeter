package service

import (
	"github.com/rodrigo462003/FlickMeter/model"
	"github.com/rodrigo462003/FlickMeter/movieAPI"
)

type MovieService interface {
	GetMovie(movieID string) (movie *model.Movie, err error)
}

type movieService struct {
	movieAPI.MovieGetter
}

func NewMovieService(apiToken string) *movieService {
	movieGetter := movieAPI.NewMovieGet(apiToken)
	return &movieService{movieGetter}
}

func (s *movieService) GetMovie(movieID string) (movie *model.Movie, err error) {
	return s.MovieGetter.GetMovie(movieID)
}
