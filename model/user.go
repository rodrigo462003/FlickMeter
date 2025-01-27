package model

import (
	"fmt"
	"net/http"
	"net/mail"
	"unicode/utf8"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(15);unique;not null"`
	Email    string `gorm:"type:varchar(254);unique;not null"`
	Password string `gorm:"type:varchar(255);not null"`
}

type UserStore interface {
	UserNameExists(string) (bool, error)
	Create(*User) error
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

func (ue userErrors) Errors() userErrors {
	return ue
}

func (ue userErrors) StatusCode() int {
	return ue.statusCode
}

func NewUser(username, email, password, confirm string, us UserStore) UserErrors {
	err := validUser(username, email, password, confirm, us)
	if err != nil {
		return err
	}

	user := &User{Username: username, Email: email, Password: password}

	if dbErr := us.Create(user); dbErr != nil {
		errors := userErrors{}
		errors.statusCode = http.StatusInternalServerError
		return errors
	}

	return nil
}

func validUser(username, email, password, confirm string, us UserStore) UserErrors {
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

	err, isDiff := ValidPassword(password, confirm)
	if err != nil {
		if isDiff {
			ue.Confirm = err.Error()
			ue.statusCode = http.StatusUnprocessableEntity
		} else {
			ue.Password = err.Error()
			ue.statusCode = http.StatusUnprocessableEntity
		}
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

func ValidPassword(p string, c string) (ValidationError, bool) {
	const (
		MaxLen = 128
		MinLen = 8
	)

	n := utf8.RuneCountInString(p)
	if n == 0 {
		return newValidationError(http.StatusUnprocessableEntity, "* Password is required."), false
	}
	if n > MaxLen {
		return newValidationErrorf(http.StatusUnprocessableEntity, "* Password must contain at most %d characters.", MaxLen), false
	}
	if n < MinLen {
		return newValidationErrorf(http.StatusUnprocessableEntity, "* Password must contain atleast %d characters.", MinLen), false
	}

	if p != c && len(c) > 0 {
		return newValidationError(http.StatusUnprocessableEntity, "* Passwords don't match."), true
	}

	return nil, false
}
