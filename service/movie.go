package service

import (
	"github.com/rodrigo462003/FlickMeter/model"
	"github.com/rodrigo462003/FlickMeter/movieAPI"
	"github.com/rodrigo462003/FlickMeter/store"
)

type MovieService interface {
	Get(movieID string) (movie *model.Movie, err error)
	Search(query string) (movies []model.Movie, err error)
	CreateReview(title, text string, movieID, userID uint) (*model.Review, error)
}

type movieService struct {
	fetcher     movieAPI.MovieFetcher
	reviewStore store.ReviewStore
}

func NewMovieService(token string, store store.ReviewStore) *movieService {
	return &movieService{movieAPI.NewMovieGet(token), store}
}

func (s *movieService) Get(movieID string) (movie *model.Movie, err error) {
	movie, err = s.fetcher.Get(movieID)
	if err != nil {
		return nil, err
	}

	reviews, err := s.reviewStore.GetByMovieID(uint(movie.ID))
	if err != nil {
		return nil, err
	}

	movie.Reviews = reviews
	return movie, nil
}

func (s *movieService) Search(query string) (movies []model.Movie, err error) {
	movies, err = s.fetcher.Search(query)
	return movies, err
}

func (s *movieService) CreateReview(title, text string, userID, movieID uint) (*model.Review, error) {
	review := model.NewReview(title, text, userID, movieID)
	if err := s.reviewStore.Create(review); err != nil {
		return nil, err
	}

	return review, nil
}
