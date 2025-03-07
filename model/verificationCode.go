package model

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"gorm.io/gorm"
)

type VerificationCode struct {
	gorm.Model
	Email     string    `gorm:"not null;index;"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

func NewVerificationCode(email string) *VerificationCode {
	num, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		panic(err)
	}

	code := fmt.Sprintf("%06d", num)
	return &VerificationCode{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
}
