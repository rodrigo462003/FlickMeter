package store

import (
	"github.com/rodrigo462003/FlickMeter/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReviewStore interface {
	Update(*model.Review) error
	GetByMovieID(movieID uint) (reviews []model.Review, err error)
	GetReview(movieID, userID uint) (review *model.Review, err error)
}

type reviewStore struct {
	db *gorm.DB
}

func NewReviewStore(db *gorm.DB) *reviewStore {
	return &reviewStore{
		db: db,
	}
}

func (rs *reviewStore) Update(review *model.Review) error {
	if err := rs.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "movie_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "text", "rating"}),
	}).Create(review).Error; err != nil {
		return err
	}

	return rs.db.Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Omit("password", "email")
	}).First(&review).Error
}

func (rs *reviewStore) GetByMovieID(movieID uint) (reviews []model.Review, err error) {
	err = rs.db.Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Omit("password", "email")
	}).Where("movie_id = ?", movieID).Find(&reviews).Error

	return reviews, err
}

func (rs *reviewStore) GetReview(movieID, userID uint) (review *model.Review, err error) {
	err = rs.db.Where("movie_id = ? AND user_id = ?", movieID, userID).First(&review).Error
	if err == nil {
		return review, nil
	}
	if err == gorm.ErrRecordNotFound {
		return &model.Review{}, nil
	}

	return nil, err
}
