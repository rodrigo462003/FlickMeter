package model

import (
	"errors"

	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	MovieID uint   `gorm:"not null;index;uniqueIndex:idx_movie_user"`
	UserID  uint   `gorm:"not null;index;uniqueIndex:idx_movie_user"`
	Title   string `gorm:"type:text;not null"`
	Text    string `gorm:"type:text;not null"`
	Rating  uint   `gorm:"not null"`

	User User `gorm:"foreignKey:UserID"`
}

func NewReview(title, text string, rating, movieID, userID uint) *Review {
	return &Review{MovieID: movieID, UserID: userID, Title: title, Text: text, Rating: rating}
}

func (r *Review) Validate() error {
	if r.Rating <= 0 || r.Rating > 10 {
		return errors.New("Invalid Rating, must be in bigger than 0 and at most to 10.")
	}

	return nil
}
