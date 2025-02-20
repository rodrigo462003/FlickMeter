package model

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/mail"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/rivo/uniseg"
	"github.com/rodrigo462003/FlickMeter/hashing"
	"gorm.io/gorm"
)

type VerificationCode struct {
	gorm.Model
	UserID    uint      `gorm:"not null;index;"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

type User struct {
	gorm.Model
	Username          string             `gorm:"type:varchar(15);unique;not null"`
	Email             string             `gorm:"type:varchar(254);unique;not null"`
	Password          string             `gorm:"type:varchar(255);not null"`
	Verified          bool               `gorm:"default:false;not null"`
	VerificationCodes []VerificationCode `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func ValidateUsername(username string) error {
	const (
		maxLen = 15
		minLen = 3
	)

	n := 0
	for _, r := range username {
		n++
		if n > maxLen {
			break
		}
		isLetter := 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z'
		if isLetter {
			continue
		}

		isDigit := '0' <= r && r <= '9'
		if isDigit {
			continue
		}

		if r == '_' || r == '-' {
			continue
		}
		return errors.New("* English letters, digits, _ and - only.")
	}

	if n == 0 {
		return errors.New("* Username is required.")
	}
	if n > maxLen {
		return fmt.Errorf("* Username must have at most %d characters.", maxLen)
	}
	if n < minLen {
		return fmt.Errorf("* Username must have at least %d characters.", minLen)
	}

	return nil
}

func ValidateEmail(email string) error {
	if len(email) == 0 {
		return errors.New("* Email address is required.")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("* This is not a valid email address.")
	}

	return nil
}

func ValidatePassword(password string) error {
	const (
		MaxLen = 128
		MinLen = 8
	)

	if len(password) == 0 {
		return errors.New("* Password is required.")
	}

	if !utf8.ValidString(password) {
		return errors.New("* Invalid character(s) detected, try again.")
	}

	n := uniseg.GraphemeClusterCount(password)
	if n > MaxLen {
		return fmt.Errorf("* Must contain at most %d characters.", MaxLen)
	}
	if n < MinLen {
		return fmt.Errorf("* Must contain at least %d characters.", MinLen)
	}

	return nil
}

func NewUser(username, email, password string) *User {
	return &User{Username: username, Email: email, Password: password}
}

func (u *User) HashPassword() error {
	hash, err := hashing.HashPassword([]byte(u.Password))
	if err != nil {
		return err
	}

	u.Password = hash
	return nil
}

func (u *User) RemoveExpiredCodes() {
	now := time.Now()
	for i := len(u.VerificationCodes) - 1; i >= 0; i-- {
		if u.VerificationCodes[i].ExpiresAt.Before(now) {
			u.VerificationCodes = append(u.VerificationCodes[:i], u.VerificationCodes[i+1:]...)
		}
	}
}

func (u *User) NewVerificationCode() error {
	code := ""
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(9))
		if err != nil {
			return err
		}
		code += strconv.FormatUint(num.Uint64(), 10)
	}

	vc := VerificationCode{
		UserID:    u.ID,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	u.VerificationCodes = append(u.VerificationCodes, vc)
	return nil
}
