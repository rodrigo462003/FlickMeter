package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Auth struct {
	gorm.Model
	UUID      uuid.UUID `gorm:"not null;index;type:uuid;"`
	UserID    uint      `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type Session struct {
	UUID      uuid.UUID
	UserID    uint
	ExpiresIn time.Duration
}

func NewSession(userID uint) *Session {
	const Day = time.Hour * 24

	return &Session{
		UUID:      uuid.New(),
		UserID:    userID,
		ExpiresIn: Day,
	}
}

func NewAuth(userID uint) *Auth {
	const Month = time.Hour * 24 * 30

	return &Auth{
		UUID:      uuid.New(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(Month),
	}
}
