package model

import (
	"fmt"
	"net/mail"
	"unicode/utf8"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username username `gorm:"type:varchar(15);unique;not null"`
	Email    email    `gorm:"type:varchar(254);unique;not null"`
	password hash     `gorm:"type:varchar(255);not null"`
}

func NewUser(u ValidUsername, e ValidEmail, p ValidPassword) *User {
	return &User{Username: u.username(), Email: e.email(), password: p.hash()}
}

type username string
type ValidUsername interface {
	String() string
	username() username
}

func (u username) username() username {
	return u
}

func (u username) String() string {
	return string(u)
}

func NewUsername(u string) (ValidUsername, error) {
	const (
		maxLen = 15
		minLen = 3
	)

	n := utf8.RuneCountInString(u)
	if n == 0 {
		return nil, fmt.Errorf("* Username is required.")
	}
	if n > maxLen {
		return nil, fmt.Errorf("* Username must have at most %d characters.", maxLen)
	}
	if n < minLen {
		return nil, fmt.Errorf("* Username must have at least %d characters.", minLen)
	}

	for _, r := range u {
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

		return nil, fmt.Errorf("* English letters, digits, _ and - only.")
	}

	return username(u), nil
}

type email string
type ValidEmail interface {
	email() email
}

func (e email) email() email {
	return e
}

func NewEmail(e string) (ValidEmail, error) {
	if len(e) == 0 {
		return nil, fmt.Errorf("* Email address is required.")
	}

	if _, err := mail.ParseAddress(string(e)); err != nil {
		return nil, fmt.Errorf("* This is not a valid email address.")
	}

	return email(e), nil
}

type hash string
type ValidPassword interface {
	hash() hash
}

func (p hash) hash() hash {
	return p
}

func NewPassword(p string) (ValidPassword, error) {
	const (
		MaxLen = 128
		MinLen = 8
	)

	n := utf8.RuneCountInString(p)
	if n == 0 {
		return nil, fmt.Errorf("* Password is required.")
	}
	if n > MaxLen {
		return nil, fmt.Errorf("* Password must contain at most %d characters.", MaxLen)
	}
	if n < MinLen {
		return nil, fmt.Errorf("* Password must contain atleast %d characters.", MinLen)
	}

	return hash(p), nil
}
