package service

import (
	"github.com/rodrigo462003/FlickMeter/model"
	"github.com/rodrigo462003/FlickMeter/movieAPI"
)

type MovieService interface {
	Get(movieID string) (movie *model.Movie, err error)
	Search(query string) (movies []model.Movie, err error)
}

type movieService struct {
	fetcher movieAPI.MovieFetcher
}

func NewMovieService(token string) *movieService {
	return &movieService{movieAPI.NewMovieGet(token)}
}

func (s *movieService) Get(movieID string) (movie *model.Movie, err error) {
	movie, err = s.fetcher.Get(movieID)
	if err != nil {
		return nil, err
	}

	videos, err := s.fetcher.Videos(movieID)
	if err != nil {
		return nil, err
	}
	movie.Videos = videos

	return movie, nil
}

func (s *movieService) Search(query string) (movies []model.Movie, err error) {
	movies, err = s.fetcher.Search(query)
	return movies[:min(5, len(movies))], err
}
