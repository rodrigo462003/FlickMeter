package service

import (
	"github.com/rodrigo462003/FlickMeter/model"
	"github.com/rodrigo462003/FlickMeter/movieAPI"
	"github.com/rodrigo462003/FlickMeter/store"
)

type MovieService interface {
	Get(movieID uint) (movie *model.Movie, err error)
	Search(query string) (movies []model.Movie, err error)
	CreateReview(title, text string, rating, movieID, userID uint) (*model.Review, error)
	GetReview(movieID uint, userID uint) (*model.Review, error)
}

type movieService struct {
	fetcher     movieAPI.MovieFetcher
	reviewStore store.ReviewStore
}

func NewMovieService(token string, store store.ReviewStore) *movieService {
	return &movieService{movieAPI.NewMovieGet(token), store}
}

func (s *movieService) GetReview(movieID, userID uint) (*model.Review, error) {
	return s.reviewStore.GetReview(movieID, userID)
}

func (s *movieService) Get(movieID uint) (movie *model.Movie, err error) {
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

func (s *movieService) CreateReview(title, text string, rating, userID, movieID uint) (*model.Review, error) {
	review := model.NewReview(title, text, rating, userID, movieID)
	if err := review.Validate(); err != nil {
		return nil, NewValidationError(err.Error(), ErrUnprocessable)
	}

	if err := s.reviewStore.Create(review); err != nil {
		return nil, err
	}

	return review, nil
}
