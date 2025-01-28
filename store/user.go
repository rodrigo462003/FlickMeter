package store

import (
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
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

func (us *UserStore) Create(user *model.User) *model.CreateUserError {
	var pgErr *pgconn.PgError
	if err := us.db.Create(user).Error; err != nil {
		cUserErrors := model.CreateUserError{}
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "uni_users_email":
				log.Println("USER EMAIL ERROR")
				cUserErrors.Email = err
			case "uni_users_username":
				cUserErrors.Username = err
			}
			return &cUserErrors
		}
		cUserErrors.Other = err
		return &cUserErrors
	}
	return nil
}
