package store

import (
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

func (us *UserStore) UsernameExists(username string) (exists bool, err error) {
	err = us.db.Raw("SELECT EXISTS (SELECT 1 FROM users WHERE LOWER(username) = LOWER(?) AND verified = TRUE)", username).Scan(&exists).Error
	return exists, err
}

func (us *UserStore) GetByID(id uint) (user *model.User, err error) {
	err = us.db.Preload(clause.Associations).First(&user, id).Error
	return user, err
}

func (us *UserStore) GetByEmail(email string) (user *model.User, err error) {
	err = us.db.Preload(clause.Associations).Where("email = ?", email).First(&user).Error
	return user, err
}

func (us *UserStore) FirstOrCreate(user *model.User) (bool, error) {
	result := us.db.Preload("VerificationCodes").Where("username = ? OR email = ?", user.Username, user.Email).FirstOrCreate(&user)
	if err := result.Error; err != nil {
		return false, err
	}

	if result.RowsAffected != 0 {
		return true, nil
	}

	return false, nil
}

func (us *UserStore) Save(u *model.User) error {
	return us.db.Save(u).Error
}

func (us *UserStore) CreateVCode(vc *model.VerificationCode) error {
	return us.db.Create(vc).Error
}

func (us *UserStore) DeleteVCode(vc *model.VerificationCode) error {
	return us.db.Delete(vc).Error
}
