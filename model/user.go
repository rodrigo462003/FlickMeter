package model

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"net/mail"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/rivo/uniseg"
	"github.com/rodrigo462003/FlickMeter/email"
	"github.com/rodrigo462003/FlickMeter/hashing"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username          string             `gorm:"type:varchar(15);unique;not null"`
	Email             string             `gorm:"type:varchar(254);unique;not null"`
	Password          string             `gorm:"type:varchar(255);not null"`
	Verified          bool               `gorm:"default:false;not null"`
	VerificationCodes []VerificationCode `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type VerificationCode struct {
	gorm.Model
	UserID    uint      `gorm:"not null;index;"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

func (u *User) hashPassword() error {
	hash, err := hashing.HashPassword([]byte(u.Password))
	if err != nil {
		return err
	}

	u.Password = hash
	return nil
}

func newVerificationCode(User *User, us UserStore) (string, StatusCoder) {
	if len(User.VerificationCodes) == 5 {
		anyRemoved := false
		for _, vc := range User.VerificationCodes {
			if vc.ExpiresAt.Before(time.Now()) {
				if err := us.DeleteVCode(&vc); err == nil {
					anyRemoved = true
				}
			}
		}
		if !anyRemoved {
			return "", NewStatusError(http.StatusConflict, "Can't request more codes. Try again later.")
		}
	}

	code := ""
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(9))
		if err != nil {
			slog.Error(err.Error())
			return "", InternalServerError()
		}
		code += strconv.FormatUint(num.Uint64(), 10)
	}

	vc := &VerificationCode{
		UserID:    User.ID,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	err := us.CreateVCode(vc)
	if err != nil {
		slog.Error(err.Error())
		return "", InternalServerError()
	}

	return code, nil
}

func NewUser(username, email, password string) *User {
	return &User{Username: username, Email: email, Password: password}
}

func NewUsername(username string) error {
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

func NewEmail(email string) error {
	if len(email) == 0 {
		return errors.New("* Email address is required.")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("* This is not a valid email address.")
	}

	return nil
}

func NewPassword(password string) error {
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
