package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rodrigo462003/FlickMeter/model"
	"gorm.io/gorm"
)

type SessionStore interface {
	Create(*model.Session) error
	Renew(*model.Session) error
}

type sessionStore struct {
	redis *redis.Client
	db    *gorm.DB
}

func NewSessionStore(redisAddr string, db *gorm.DB) *sessionStore {
	redis := redis.NewClient(&redis.Options{Addr: redisAddr})
	return &sessionStore{redis, db}
}

func (store *sessionStore) Create(s *model.Session) error {
	if s.Persist {
		return store.db.Create(s).Error
	}
	return store.redis.Set(context.Background(), s.UUID.String(), s.UserID, s.ExpiresIn).Err()
}

func (store *sessionStore) Renew(s *model.Session) error {
	return store.redis.Expire(context.Background(), s.UUID.String(), time.Hour*24).Err()
}
