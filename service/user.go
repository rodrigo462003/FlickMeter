package service

import (
	"fmt"
	"slices"
	"time"

	"github.com/rodrigo462003/FlickMeter/email"
	"github.com/rodrigo462003/FlickMeter/model"
	"github.com/rodrigo462003/FlickMeter/store"
)

type UserService interface {
	ValidatePassword(password string) ValidationError
	ValidateEmail(email string) ValidationError
	ValidateUsername(username string) error
	Register(username, email, password string) error
	Verify(code, username, email, password string) (*model.Session, error)
	SignIn(email, password string) (uint, error)
	CreateSession(userID uint) (*model.Session, error)
	CreateAuth(userID uint) (*model.Auth, error)
	GetUserFromAuth(uuid string) (*model.User, error)
	GetUserFromSession(uuid string) (*model.User, error)
	DeleteSession(uuid string) error
	DeleteAuth(uuid string) error
}

type userService struct {
	userStore    store.UserStore
	sender       email.EmailSender
	sessionStore store.SessionStore
}

func NewUserService(us store.UserStore, ss store.SessionStore, es email.EmailSender) *userService {
	return &userService{us, es, ss}
}

func (s *userService) DeleteSession(uuid string) error {
	return s.sessionStore.DeleteSession(uuid)
}

func (s *userService) DeleteAuth(uuid string) error {
	return s.sessionStore.DeleteAuth(uuid)
}

func (s *userService) CreateSession(userID uint) (*model.Session, error) {
	session := model.NewSession(userID)
	if err := s.sessionStore.CreateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *userService) CreateAuth(userID uint) (*model.Auth, error) {
	auth := model.NewAuth(userID)
	if err := s.sessionStore.CreateAuth(auth); err != nil {
		return nil, err
	}

	return auth, nil
}

func (s *userService) GetUserFromSession(uuid string) (*model.User, error) {
	id, err := s.sessionStore.GetUserIDBySession(uuid)
	if err != nil {
		return nil, err
	}

	return s.userStore.GetByID(id)
}

func (s *userService) GetUserFromAuth(uuid string) (*model.User, error) {
	user, err := s.sessionStore.GetUserByAuth(uuid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) SignIn(email, password string) (uint, error) {
	user, err := s.userStore.GetByEmail(email)
	if err != nil {
		return 0, NewValidationError("* Incorrect Email or Password.", ErrUnauthorized)
	}

	if !user.PasswordsMatch(password) {
		return 0, NewValidationError("* Incorrect Email or Password.", ErrUnauthorized)
	}

	return user.ID, nil
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

	code := model.NewVerificationCode(email)

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

func (s *userService) Verify(code, username, email, password string) (*model.Session, error) {
	if err := s.isValidCode(email, code); err != nil {
		return nil, err
	}

	userID, err := s.createUser(username, email, password)
	if err != nil {
		return nil, err
	}

	session := model.NewSession(userID)
	if err := s.sessionStore.CreateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *userService) isValidCode(email, code string) error {
	vCodes, err := s.userStore.GetVCodesByEmail(email)
	if err != nil {
		return err
	}

	now := time.Now()
	if !slices.ContainsFunc(vCodes, func(dbCode model.VerificationCode) bool {
		return dbCode.ExpiresAt.After(now) && code == dbCode.Code
	}) {
		return NewValidationError("* Incorrect code, try again.", ErrConflict)
	}

	return nil
}

func (s *userService) createUser(username, email, password string) (uint, error) {
	user := model.NewUser(username, email, password)
	user.HashPassword()

	if err := s.userStore.Create(user); err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			return 0, NewValidationError("* Incorrect code, try again", ErrConflict)
		case store.ErrDuplicateUsername:
			return 0, NewValidationErrorsSingle("username", "* Username already taken.", ErrConflict)
		}
		return 0, err
	}

	return user.ID, nil
}
