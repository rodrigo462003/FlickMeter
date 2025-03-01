package service

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/rodrigo462003/FlickMeter/email"
	"github.com/rodrigo462003/FlickMeter/model"
	"github.com/rodrigo462003/FlickMeter/store"
)

type UserService interface {
	ValidatePassword(password string) ValidationError
	ValidateEmail(email string) ValidationError
	ValidateUsername(username string) error
	Register(username, email, password string) error
	Verify(code, username, email, password string) (*http.Cookie, error)
}

type userService struct {
	userStore    store.UserStore
	sender       email.EmailSender
	sessionStore store.SessionStore
}

func NewUserService(us store.UserStore, ss store.SessionStore, es email.EmailSender) *userService {
	return &userService{us, es, ss}
}

func (s *userService) ValidatePassword(password string) ValidationError {
	if err := model.ValidatePassword(password); err != nil {
		return NewValidationError(err.Error(), ErrUnprocessable)
	}

	return nil
}

func (s *userService) ValidateEmail(email string) ValidationError {
	if err := model.ValidateEmail(email); err != nil {
		return NewValidationError(err.Error(), ErrUnprocessable)
	}

	return nil
}

func (s *userService) ValidateUsername(username string) error {
	if err := model.ValidateUsername(username); err != nil {
		return NewValidationError(err.Error(), ErrUnprocessable)
	}

	isDupe, err := s.userStore.UsernameExists(username)
	if err != nil {
		return err
	}
	if isDupe {
		return NewValidationError("* Username already taken.", ErrConflict)
	}

	return nil
}

func (s *userService) validateUser(user *model.User) error {
	vErrs := NewValidationErrors()
	if err := s.ValidateUsername(user.Username); err != nil {
		if vErr, ok := err.(ValidationError); ok {
			vErrs.add("username", vErr)
		} else {
			return err
		}
	}
	if err := s.ValidateEmail(user.Email); err != nil {
		vErrs.add("email", err)
	}

	if err := s.ValidatePassword(user.Password); err != nil {
		vErrs.add("password", err)
	}

	if len(vErrs.errorMap) > 0 {
		return vErrs
	}

	return nil
}

func (s *userService) removeExpired(vCodes *[]model.VerificationCode) error {
	var err error
	now := time.Now()
	*vCodes = slices.DeleteFunc(*vCodes, func(vCode model.VerificationCode) bool {
		if !vCode.ExpiresAt.Before(now) {
			return false
		}
		if err = s.userStore.DeleteVCode(vCode); err != nil {
			return false
		}
		return true
	})

	return err
}

func (s *userService) createVerificationCode(email string) (*model.VerificationCode, error) {
	vCodes, err := s.userStore.GetVCodesByEmail(email)
	if err != nil {
		return nil, err
	}
	if len(vCodes) == 5 {
		if err := s.removeExpired(&vCodes); len(vCodes) == 5 {
			if err != nil {
				return nil, err
			}
			return nil, NewValidationError("* Can't request more codes. Try again later.", ErrConflict)
		}
	}

	code, err := model.NewVerificationCode(email)
	if err != nil {
		return nil, err
	}

	if err := s.userStore.CreateVCode(code); err != nil {
		return nil, err
	}

	return code, nil
}

func (s *userService) Register(username, email, password string) error {
	const (
		emailSubject = "FlickMeter registration"
		codeBodyF    = "Please enter the following code to complete your Signup.\r\n%s"
		verifiedBody = "You already have a FlickMeter account, try signing in."
	)

	user := model.NewUser(username, email, password)
	if err := s.validateUser(user); err != nil {
		return err
	}

	exists, err := s.userStore.EmailExists(user.Email)
	if err != nil {
		return err
	}
	if exists {
		s.sender.SendMail(user.Email, emailSubject, verifiedBody)
		return nil
	}

	code, err := s.createVerificationCode(user.Email)
	if err != nil {
		return err
	}

	s.sender.SendMail(user.Email, emailSubject, fmt.Sprintf(codeBodyF, code.Code))

	return nil
}

func (s *userService) Verify(code, username, email, password string) (*http.Cookie, error) {
	if err := s.isValidCode(email, code); err != nil {
		return nil, err
	}

	if err := s.createUser(username, email, password); err != nil {
		return nil, err
	}

	cookie, err := newCookie()
	if err != nil {
		return nil, err
	}

	return cookie, nil
}

func newCookie() (*http.Cookie, error) {
	return &http.Cookie{
		Name:     "session",
		Value:    uuid.NewString(),
		Path:     "/",
		Secure:   false, //SET TO SECURE FOR HTTPS.
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}, nil
}

func (s *userService) isValidCode(email, code string) error {
	vCodes, err := s.userStore.GetVCodesByEmail(email)
	if err != nil {
		return err
	}

	now := time.Now()
	fmt.Println(email, vCodes, code)
	if !slices.ContainsFunc(vCodes, func(dbCode model.VerificationCode) bool {
		return dbCode.ExpiresAt.After(now) && code == dbCode.Code
	}) {
		return NewValidationError("* Incorrect code, try again.", ErrConflict)
	}

	return nil
}

func (s *userService) createUser(username, email, password string) error {
	user := model.NewUser(username, email, password)
	user.MustHashPassword()

	if err := s.userStore.Create(user); err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			return NewValidationError("* Incorrect code, try again", ErrConflict)
		case store.ErrDuplicateUsername:
			return NewValidationErrorsSingle("username", "* Username already taken.", ErrConflict)
		}
		return err
	}

	return nil
}
