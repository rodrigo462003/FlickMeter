package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rodrigo462003/FlickMeter/model"
)

type SessionStore interface {
	Create(*model.Session) error
	Renew(*model.Session) error
}

type sessionStore struct {
	db *redis.Client
}

func NewSessionStore(addr string) *sessionStore {
	options := &redis.Options{Addr: addr}
	return &sessionStore{redis.NewClient(options)}
}

func (store *sessionStore) Create(s *model.Session) error {
	return store.db.Set(context.Background(), s.UUID.String(), s.UserID, s.Expire).Err()
}

func (store *sessionStore) Renew(s *model.Session) error {
	return store.db.Expire(context.Background(), s.UUID.String(), time.Hour*24).Err()
}
