package store

import (
	"github.com/rodrigo462003/FlickMeter/model"
	"gorm.io/gorm"
)

type ReviewStore interface {
	Create(*model.Review) error
	GetByMovieID(movieID uint) (reviews []model.Review, err error)
}

type reviewStore struct {
	db *gorm.DB
}

func NewReviewStore(db *gorm.DB) *reviewStore {
	return &reviewStore{
		db: db,
	}
}

func (rs *reviewStore) Create(review *model.Review) error {
	return rs.db.Create(&review).Error
}

func (rs *reviewStore) GetByMovieID(movieID uint) (reviews []model.Review, err error) {
	err = rs.db.Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Omit("password", "email")
	}).Where("movie_id = ?", movieID).Find(&reviews).Error

	return reviews, err
}
