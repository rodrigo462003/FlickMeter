package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UUID   uuid.UUID
	UserID uint
	Expire time.Duration
}

func NewSession(userID uint) *Session {
	return &Session{uuid.New(), userID, time.Hour * 24}
}
