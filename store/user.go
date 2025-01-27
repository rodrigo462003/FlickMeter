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

func (us *UserStore) UserNameExists(username string) (bool, error) {
	var exists bool
	err := us.db.Raw("SELECT EXISTS (SELECT 1 FROM users WHERE LOWER(username) = LOWER(?))", username).Scan(&exists).Error
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (us *UserStore) Create(user *model.User) error {
	return us.db.Create(user).Error
}
