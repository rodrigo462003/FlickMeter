package store

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rodrigo462003/FlickMeter/model"
	"gorm.io/gorm"
)

type UserStore interface {
	UsernameExists(string) (bool, error)
	EmailExists(string) (bool, error)
	Create(*model.User) error
	GetByID(uint) (*model.User, error)
	GetByEmail(string) (*model.User, error)
	Save(*model.User) error
	GetVCodesByEmail(string) ([]model.VerificationCode, error)
	CreateVCode(*model.VerificationCode) error
	DeleteVCode(model.VerificationCode) error
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
	err = us.db.Raw("SELECT EXISTS (SELECT 1 FROM users WHERE LOWER(username) = LOWER(?))", username).Scan(&exists).Error
	return exists, err
}

func (us *userStore) EmailExists(email string) (exists bool, err error) {
	err = us.db.Raw("SELECT EXISTS (SELECT 1 FROM users WHERE LOWER(email) = LOWER(?))", email).Scan(&exists).Error
	return exists, err
}

func (us *userStore) GetByID(id uint) (user *model.User, err error) {
	if err = us.db.Preload("Watchlist").First(&user, id).Error; err != nil {
		return nil, err
	}

	return user, err
}

func (us *userStore) GetByEmail(email string) (user *model.User, err error) {
	err = us.db.Where("email = ?", email).First(&user).Error
	return user, err
}

var ErrDuplicateEmail = errors.New("* Email already taken.")
var ErrDuplicateUsername = errors.New("* Username already taken.")

func (us *userStore) Create(user *model.User) error {
	if err := us.db.Create(&user).Error; err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			if strings.Contains(pgErr.Message, "username") {
				return ErrDuplicateUsername
			}
			if strings.Contains(pgErr.Message, "email") {
				return ErrDuplicateEmail
			}
		}

		return err
	}

	return nil
}

func (us *userStore) Save(user *model.User) error {
	return us.db.Save(user).Error
}

func (us *userStore) GetVCodesByEmail(email string) (vc []model.VerificationCode, err error) {
	err = us.db.Where("email = ?", email).Find(&vc).Error
	return vc, err
}

func (us *userStore) CreateVCode(vCode *model.VerificationCode) error {
	return us.db.Create(vCode).Error
}

func (us *userStore) DeleteVCode(vCode model.VerificationCode) error {
	return us.db.Delete(&vCode).Error
}
