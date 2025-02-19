package store

import (
	"github.com/rodrigo462003/FlickMeter/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserStore interface {
	UsernameExists(string) (bool, error)
	FirstOrCreate(*model.User) (bool, error)
	GetByID(uint) (*model.User, error)
	GetByEmail(string) (*model.User, error)
	Save(*model.User) error
}

type userStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *userStore {
	return &userStore{
		db: db,
	}
}

func (us *userStore) UsernameExists(username string) (exists bool, err error) {
	err = us.db.Raw("SELECT EXISTS (SELECT 1 FROM users WHERE LOWER(username) = LOWER(?) AND verified = TRUE)", username).Scan(&exists).Error
	return exists, err
}

func (us *userStore) GetByID(id uint) (user *model.User, err error) {
	err = us.db.Preload(clause.Associations).First(&user, id).Error
	return user, err
}

func (us *userStore) GetByEmail(email string) (user *model.User, err error) {
	err = us.db.Preload(clause.Associations).Where("email = ?", email).First(&user).Error
	return user, err
}

func (us *userStore) FirstOrCreate(user *model.User) (bool, error) {
	result := us.db.Preload("VerificationCodes").Where("username = ? OR email = ?", user.Username, user.Email).FirstOrCreate(&user)
	if err := result.Error; err != nil {
		return false, err
	}

	if result.RowsAffected != 0 {
		return true, nil
	}

	return false, nil
}

func (us *userStore) Save(user *model.User) error {
	//Make sure it saves Associations..
	return us.db.Save(&user).Error
}
