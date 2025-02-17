package store

import (
	"github.com/rodrigo462003/FlickMeter/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserStore interface {
	UsernameExists(string) (bool, error)
	FirstOrCreate(*model.User) (bool, error)
	CreateVCode(*model.VerificationCode) error
	GetByID(uint) (*model.User, error)
	GetByEmail(string) (*model.User, error)
	DeleteVCode(*model.VerificationCode) error
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

func (us *userStore) Save(u *model.User) error {
	return us.db.Save(u).Error
}

func (us *userStore) CreateVCode(vc *model.VerificationCode) error {
	return us.db.Create(vc).Error
}

func (us *userStore) DeleteVCode(vc *model.VerificationCode) error {
	return us.db.Delete(vc).Error
}
