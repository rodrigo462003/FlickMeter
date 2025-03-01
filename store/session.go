package store

import "github.com/redis/go-redis/v9"

type SessionStore interface {
}

type sessionStore struct {
	db *redis.Client
}

func NewSessionStore(addr string) *sessionStore {
	options := &redis.Options{Addr: addr}
	client := redis.NewClient(options)

	return &sessionStore{client}
}
