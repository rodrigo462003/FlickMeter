package store

import (
	"log/slog"
	"net/http"

	"github.com/rodrigo462003/FlickMeter/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (us *UserStore) GetUserByID(id uint) (user *model.User, err error) {
	err = us.db.Preload(clause.Associations).First(&user, id).Error
	return user, err
}

// Creates and stores user, returns error with StatusError for internals and email conflict.
// Return StatusErrors with conflict if username taken in map.
func (us *UserStore) Create(user *model.User) model.StatusCoder {
	tempUser := model.User{Username: user.Username, Email: user.Email}

	result := us.db.FirstOrCreate(&tempUser)
	if result.Error != nil {
		slog.Error(result.Error.Error())
		return model.NewStatusError(http.StatusInternalServerError, "Something unexpected has happened, please try again.")
	}
	if result.RowsAffected == 1 {
		*user = tempUser
		return nil
	}

	if !tempUser.Verified {
		if tempUser.Email == user.Email {
			if err := us.db.Preload("VerificationCodes").First(&tempUser).Error; err != nil {
				slog.Error(err.Error())
				return model.NewStatusError(http.StatusInternalServerError, "Something unexpected has happened, please try again.")
			}
			*user = tempUser
			return nil
		}
		errorMap := map[string]string{"username": "Username already taken."}
		return model.NewStatusErrors(http.StatusConflict, errorMap)
	}
	if tempUser.Email == user.Email {
		return model.NewStatusError(http.StatusConflict, "Email already taken.")
	}
	errorMap := map[string]string{"username": "Username already taken."}
	return model.NewStatusErrors(http.StatusConflict, errorMap)
}

func (us *UserStore) CreateVerificationCode(vc *model.VerificationCode) error {
	return us.db.Create(vc).Error
}
