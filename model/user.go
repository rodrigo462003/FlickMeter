package model

import (
	"crypto/rand"
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

type UserStore interface {
	UsernameExists(string) (bool, error)
	FirstOrCreate(*User) (bool, error)
	CreateVCode(*VerificationCode) error
	GetByID(uint) (*User, error)
	GetByEmail(string) (*User, error)
	DeleteVCode(*VerificationCode) error
	Save(*User) error
}

func (u *User) hashPassword() error {
	hash, err := hashing.HashPassword([]byte(u.Password))
	if err != nil {
		return err
	}

	u.Password = hash
	return nil
}

func VerifyUser(code string, email string, us UserStore) StatusCoder {
	user, err := us.GetByEmail(email)
	if err != nil {
		slog.Error(err.Error())
		return InternalServerError()
	}

	if user.Verified {
		return nil
	}

	for _, dbCode := range user.VerificationCodes {
		if dbCode.ExpiresAt.After(time.Now()) && code == dbCode.Code {
			user.Verified = true
			if err := us.Save(user); err != nil {
				user.Verified = false
				slog.Error(err.Error())
				return InternalServerError()
			}
			return nil
		}
	}

	return NewStatusError(http.StatusConflict, "* Incorrect code.")
}

func newVerificationCode(user *User, us UserStore) (string, StatusCoder) {
	if len(user.VerificationCodes) == 5 {
		anyRemoved := false
		for _, vc := range user.VerificationCodes {
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
		UserID:    user.ID,
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

func NewUser(username, email, password string, us UserStore, es email.EmailSender) StatusCoder {
	const emailSubject = "FlickMeter registration"

	user := &User{Username: username, Email: email, Password: password}

	if multiErr := user.isValid(us); multiErr != nil {
		return multiErr
	}

	if err := user.hashPassword(); err != nil {
		slog.Error(err.Error())
		return InternalServerError()
	}

	tmpU := *user
	created, err := us.FirstOrCreate(user)
	if err != nil {
		slog.Error(err.Error())
		return InternalServerError()
	}
	if !created {
		if user.Verified {
			if user.Username == tmpU.Username && user.Email != tmpU.Email {
				return NewStatusErrors(http.StatusConflict, map[string]string{"username": "Username already taken."})
			}
			if user.Email == tmpU.Email {
				es.SendMail(user.Email, emailSubject, "You already have a FlickMeter account, try signing in.")
				return nil
			}
			slog.Error(err.Error())
			return InternalServerError()
		}

		if user.Username == tmpU.Username && user.Email != tmpU.Email {
			return NewStatusErrors(http.StatusConflict, map[string]string{"username": "Username already taken."})
		}
		if user.Email == tmpU.Email {
			user.Username, user.Password = user.Username, user.Password
			if err := us.Save(user); err != nil {
				slog.Error(err.Error())
				return InternalServerError()
			}
		} else {
			slog.Error(err.Error())
			return InternalServerError()
		}
	}

	code, codeErr := newVerificationCode(user, us)
	if err != nil {
		return codeErr
	}

	emailBody := fmt.Sprintf("Please enter the following code to complete your Signup.\r\n%s", code)
	es.SendMail(user.Email, emailSubject, emailBody)

	return nil
}

type mutliStatusCoder interface {
	StatusCoder
	Map() map[string]string
}

type StatusErrors struct {
	code     int
	errorMap map[string]string
}

func (e StatusErrors) Error() string {
	return fmt.Sprint(e.errorMap)
}

func (e StatusErrors) StatusCode() int {
	return e.code
}

func (e StatusErrors) Map() map[string]string {
	return e.errorMap
}

func NewStatusErrors(code int, m map[string]string) StatusErrors {
	return StatusErrors{code, m}
}

func (user *User) isValid(us UserStore) StatusCoder {
	errorMap := make(map[string]string, 3)
	codes := make([]int, 0, len(errorMap))

	if err := ValidUsername(user.Username, us); err != nil {
		errorMap["username"] = err.Error()
		codes = append(codes, err.StatusCode())
	}

	if err := ValidEmail(user.Email); err != nil {
		errorMap["email"] = err.Error()
		codes = append(codes, err.StatusCode())
	}

	if err := ValidPassword(user.Password); err != nil {
		errorMap["password"] = err.Error()
		codes = append(codes, err.StatusCode())
	}

	if len(errorMap) == 0 {
		return nil
	}

	statusCode := getPriorityStatusCode(codes)

	return NewStatusErrors(statusCode, errorMap)
}

func getPriorityStatusCode(codes []int) int {
	priority := [...]int{
		http.StatusInternalServerError,
		http.StatusConflict,
		http.StatusUnprocessableEntity,
	}

	statusSet := make(map[int]struct{}, len(priority))
	for _, code := range codes {
		statusSet[code] = struct{}{}
	}

	statusCode := http.StatusInternalServerError
	for _, code := range priority {
		if _, ok := statusSet[code]; ok {
			return code
		}
	}

	return statusCode
}

func ValidUsername(u string, us UserStore) StatusCoder {
	const (
		maxLen = 15
		minLen = 3
	)

	n := 0
	for _, r := range u {
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
		return NewStatusError(http.StatusUnprocessableEntity, "* English letters, digits, _ and - only.")
	}

	if n == 0 {
		return NewStatusError(http.StatusUnprocessableEntity, "* Username is required.")
	}
	if n > maxLen {
		return NewStatusErrorf(http.StatusUnprocessableEntity, "* Username must have at most %d characters.", maxLen)
	}
	if n < minLen {
		return NewStatusErrorf(http.StatusUnprocessableEntity, "* Username must have at least %d characters.", minLen)
	}

	alreadyExists, err := us.UsernameExists(u)
	if err != nil {
		return NewStatusError(http.StatusInternalServerError, "")
	}
	if alreadyExists {
		return NewStatusError(http.StatusConflict, "* Username already taken.")
	}

	return nil
}

func ValidEmail(e string) StatusCoder {
	if len(e) == 0 {
		return NewStatusError(http.StatusUnprocessableEntity, "* Email address is required.")
	}

	if _, err := mail.ParseAddress(e); err != nil {
		return NewStatusError(http.StatusUnprocessableEntity, "* This is not a valid email address.")
	}

	return nil
}

func ValidPassword(p string) StatusCoder {
	const (
		MaxLen = 128
		MinLen = 8
	)

	if !utf8.ValidString(p) {
		return NewStatusError(http.StatusUnprocessableEntity, "* Invalid character(s) detected, try again.")
	}

	n := uniseg.GraphemeClusterCount(p)
	if n == 0 {
		return NewStatusError(http.StatusUnprocessableEntity, "* Password is required.")
	}
	if n > MaxLen {
		return NewStatusErrorf(http.StatusUnprocessableEntity, "* Must contain at most %d characters.", MaxLen)
	}
	if n < MinLen {
		return NewStatusErrorf(http.StatusUnprocessableEntity, "* Must contain at least %d characters.", MinLen)
	}

	return nil
}

type StatusCoder interface {
	error
	StatusCode() int
}

type StatusError struct {
	code    int
	message string
}

func NewStatusError(code int, message string) StatusError {
	return StatusError{code, message}
}

func InternalServerError() StatusError {
	return NewStatusError(http.StatusInternalServerError, "Something unexpected has happened, please try again.")
}

func NewStatusErrorf(code int, format string, a ...any) StatusError {
	message := fmt.Sprintf(format, a...)
	return StatusError{code, message}
}

func (e StatusError) Error() string {
	return e.message
}

func (e StatusError) StatusCode() int {
	return e.code
}
