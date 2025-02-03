package store

import (
	"errors"
	"fmt"

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

func (us *UserStore) Create(user *model.User) (uint, *model.CreateUserError) {
	tempUser := model.User{Username: user.Username, Email: user.Email}

	result := us.db.FirstOrCreate(&tempUser)
	if result.Error != nil {
		return 0, &model.CreateUserError{Other: errors.New("Failed to FirstOrCreate(User)")}
	}
	if result.RowsAffected == 1 {
		return tempUser.ID, nil
	}
	fmt.Println(user, tempUser)

	if !tempUser.Verified {
		if tempUser.Email == user.Email {
			return tempUser.ID, nil
		}
		if tempUser.Username == user.Username {
			return 0, &model.CreateUserError{Username: errors.New("Username already exists.")}
		}
	} else {
		if tempUser.Username == user.Username {
			return 0, &model.CreateUserError{Username: errors.New("Username already exists.")}
		}
		if tempUser.Email == user.Email {
			return 0, &model.CreateUserError{Email: errors.New("Email already exits.")}
		}
	}
	return 0, &model.CreateUserError{Other: errors.New("Error creating user: this shouldnt be possible.")}
}

func (us *UserStore) CreateVerificationCode(vc *model.VerificationCode) error {
	return us.db.Create(vc).Error
}
