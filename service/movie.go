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
