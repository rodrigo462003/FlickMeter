package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rodrigo462003/FlickMeter/model"
	"gorm.io/gorm"
)

type SessionStore interface {
	CreateSession(session *model.Session) error
	CreateAuth(auth *model.Auth) error
	GetUserIDBySession(uuid string) (uint, error)
	GetUserByAuth(uuid string) (*model.User, error)
}

type sessionStore struct {
	redis *redis.Client
	db    *gorm.DB
}

func NewSessionStore(redisAddr string, db *gorm.DB) *sessionStore {
	redis := redis.NewClient(&redis.Options{Addr: redisAddr})
	return &sessionStore{redis, db}
}

func (store *sessionStore) CreateSession(s *model.Session) error {
	return store.redis.Set(context.TODO(), s.UUID.String(), s.UserID, s.ExpiresIn).Err()
}

func (store *sessionStore) CreateAuth(s *model.Auth) error {
	return store.db.Create(s).Error
}

func (store *sessionStore) GetUserByAuth(uuid string) (*model.User, error) {
	tx := store.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var auth model.Auth
	if err := tx.Preload("User").Where("uuid = ?", uuid).First(&auth).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	auth.ExpiresAt = time.Now().Add(30 * 24 * time.Hour)
	if err := tx.Save(&auth).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &auth.User, nil
}

func (store *sessionStore) GetUserIDBySession(uuid string) (uint, error) {
	pipe := store.redis.TxPipeline()

	get := pipe.Get(context.TODO(), uuid)
	pipe.Expire(context.TODO(), uuid, 24*time.Hour)

	_, err := pipe.Exec(context.TODO())
	if err != nil {
		return 0, err
	}

	val, err := get.Uint64()
	if err != nil {
		return 0, err
	}

	return uint(val), nil
}
