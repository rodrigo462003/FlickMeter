package store

import (
	"github.com/rodrigo462003/FlickMeter/model"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us UserStore) UserNameExists(username model.ValidUsername) (bool, error) {
	var exists bool
	err := us.db.Raw("SELECT EXISTS (SELECT 1 FROM users WHERE LOWER(username) = LOWER(?))", username.String()).Scan(&exists).Error
	if err != nil {
		return false, err
	}

	return exists, nil
}
