package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	UUID      uuid.UUID     `gorm:"not null;index;type:uuid;"`
	UserID    uint          `gorm:"not null"`
	ExpiresAt time.Time     `gorm:"not null"`
	User      User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	ExpiresIn time.Duration `gorm:"-"`
	Persist   bool          `gorm:"-"`
}

func NewSession(userID uint, persist bool) *Session {
	const Day = time.Hour * 24
	const Month = time.Hour * 24 * 30

	s := &Session{
		UUID:    uuid.New(),
		UserID:  userID,
		Persist: persist,
	}

	if persist {
		s.ExpiresAt = time.Now().Add(Month)
	} else {
		s.ExpiresIn = Day
	}

	return s
}
