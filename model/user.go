package model

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"net/mail"
	"time"
	"unicode/utf8"

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
	VerificationCodes []VerificationCode `gorm:"foreignKey:UserID"`
}

type VerificationCode struct {
	gorm.Model
	UserID    uint      `gorm:"not null;index"`
	Code      uint      `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

type UserStore interface {
	UserNameExists(string) (bool, error)
	Create(*User) (uint, *CreateUserError)
	CreateVerificationCode(*VerificationCode) error
}

type UserErrors interface {
	Errors() userErrors
	StatusCode() int
}

type userErrors struct {
	Username   string
	Email      string
	Password   string
	Confirm    string
	statusCode int
}

type CreateUserError struct {
	Username error
	Email    error
	Other    error
}

func (ue userErrors) Errors() userErrors {
	return ue
}

func (ue userErrors) StatusCode() int {
	return ue.statusCode
}

func (u *User) hashPassword() error {
	hash, err := hashing.HashPassword([]byte(u.Password))
	if err != nil {
		return err
	}

	u.Password = hash
	return nil
}

func NewVerificationCode(userID uint, us UserStore) (uint, error) {
	code := uint(0)
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(9))
		if err != nil {
			return 0, err
		}
		code = code*10 + uint(num.Uint64()) + 1
	}

	vc := &VerificationCode{
		UserID:    userID,
		Code:      uint(code),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	err := us.CreateVerificationCode(vc)
	if err != nil {
		return 0, err
	}

	return code, nil
}

func NewUser(username, email, password string, us UserStore, es email.EmailSender) UserErrors {
	const emailSubject = "FlickMeter registration"

	errI := validUser(username, email, password, us)
	if errI != nil {
		return errI
	}

	user := &User{Username: username, Email: email, Password: password}

	if hashErr := user.hashPassword(); hashErr != nil {
		return userErrors{statusCode: http.StatusInternalServerError}
	}

	emailBody := "You already have an account, try logging in."
	id, dbErr := us.Create(user)
	if dbErr != nil {
		if dbErr.Email != nil {
		} else if dbErr.Username != nil {
			return userErrors{Username: "Username is already taken.", statusCode: http.StatusConflict}
		} else {
			return userErrors{statusCode: http.StatusInternalServerError}
		}
	}

	code, codeErr := NewVerificationCode(id, us)
	if codeErr != nil {
		return userErrors{statusCode: http.StatusInternalServerError}
	}

	emailBody = fmt.Sprintf("Please enter the following code to complete your Signup.\r\n%d", code)

	es.SendMail(user.Email, emailSubject, emailBody)

	return nil
}

func validUser(username, email, password string, us UserStore) UserErrors {
	ue := userErrors{}
	ue.statusCode = http.StatusOK

	err := ValidUsername(username, us)
	if err != nil {
		if err.StatusCode() == http.StatusInternalServerError {
			ue.statusCode = err.StatusCode()
			return ue
		}
		ue.Username = err.Error()
		ue.statusCode = err.StatusCode()
	}

	err = ValidEmail(email)
	if err != nil {
		ue.Email = err.Error()
		ue.statusCode = http.StatusUnprocessableEntity
	}

	err = ValidPassword(password)
	if err != nil {
		ue.Password = err.Error()
		ue.statusCode = http.StatusUnprocessableEntity
	}
	if ue.statusCode != http.StatusOK {
		return ue
	}
	return nil
}

type ValidationError interface {
	error
	StatusCode() int
}

type validationError struct {
	message    string
	statusCode int
}

func newValidationError(statusCode int, message string) validationError {
	return validationError{message, statusCode}
}

func newValidationErrorf(statusCode int, format string, args ...any) validationError {
	message := fmt.Sprintf(format, args...)
	return validationError{
		message:    message,
		statusCode: statusCode,
	}
}

func (v validationError) Error() string {
	return v.message
}

func (v validationError) StatusCode() int {
	return v.statusCode
}

func ValidUsername(u string, us UserStore) ValidationError {
	const (
		maxLen = 15
		minLen = 3
	)

	n := utf8.RuneCountInString(u)
	if n == 0 {
		return newValidationError(http.StatusUnprocessableEntity, "* Username is required.")
	}
	if n > maxLen {
		return newValidationErrorf(http.StatusUnprocessableEntity, "* Username must have at most %d characters.", maxLen)
	}
	if n < minLen {
		return newValidationErrorf(http.StatusUnprocessableEntity, "* Username must have at least %d characters.", minLen)
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
		return newValidationError(http.StatusUnprocessableEntity, "* English letters, digits, _ and - only.")
	}

	alreadyExists, err := us.UserNameExists(u)
	if err != nil {
		return newValidationError(http.StatusInternalServerError, err.Error())
	}
	if alreadyExists {
		return newValidationError(http.StatusConflict, "Username already taken.")
	}

	return nil
}

func ValidEmail(e string) ValidationError {
	if len(e) == 0 {
		return newValidationError(http.StatusUnprocessableEntity, "* Email address is required.")
	}

	if _, err := mail.ParseAddress(string(e)); err != nil {
		return newValidationError(http.StatusUnprocessableEntity, "* This is not a valid email address.")
	}

	return nil
}

func ValidPassword(p string) ValidationError {
	const (
		MaxLen = 128
		MinLen = 8
	)

	n := utf8.RuneCountInString(p)
	if n == 0 {
		return newValidationError(http.StatusUnprocessableEntity, "* Password is required.")
	}
	if n > MaxLen {
		return newValidationErrorf(http.StatusUnprocessableEntity, "* Password must contain at most %d characters.", MaxLen)
	}
	if n < MinLen {
		return newValidationErrorf(http.StatusUnprocessableEntity, "* Password must contain atleast %d characters.", MinLen)
	}

	return nil
}
