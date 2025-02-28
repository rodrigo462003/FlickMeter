package model

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type VerificationCode struct {
	gorm.Model
	Email     string    `gorm:"not null;index;"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

func NewVerificationCode(email string) (*VerificationCode, error) {
	code := ""
	for range 6 {
		num, err := rand.Int(rand.Reader, big.NewInt(9))
		if err != nil {
			return nil, err
		}
		code += strconv.FormatUint(num.Uint64(), 10)
	}

	return &VerificationCode{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}, nil
}
