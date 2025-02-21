package model

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(15);unique;not null"`
	Email    string `gorm:"type:varchar(254);unique;not null"`
	Password string `gorm:"type:varchar(255);not null"`
}

func NewVerificationCode(email string) (*VerificationCode, error) {
	code := ""
	for i := 0; i < 6; i++ {
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
