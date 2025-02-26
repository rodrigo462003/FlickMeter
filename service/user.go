package service

import (
	"fmt"
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
	Verify(code, username, email, password string) error
}

type userService struct {
	store  store.UserStore
	sender email.EmailSender
}

func NewUserService(us store.UserStore, es email.EmailSender) *userService {
	return &userService{us, es}
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

	isDupe, err := s.store.UsernameExists(username)
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

func (s *userService) removeExpired(verificationCodes []model.VerificationCode) ([]model.VerificationCode, error) {
	expiredCodes := make([]model.VerificationCode, 0, 5)
	nonExpiredCodes := make([]model.VerificationCode, 0, 5)

	for _, vCode := range verificationCodes {
		if vCode.ExpiresAt.Before(time.Now()) {
			expiredCodes = append(expiredCodes, vCode)
		} else {
			nonExpiredCodes = append(nonExpiredCodes, vCode)
		}
	}

	if len(expiredCodes) > 0 {
		if err := s.store.DeleteVCodes(expiredCodes); err != nil {
			return nil, err
		}
	}

	return nonExpiredCodes, nil
}

func (s *userService) createVerificationCode(email string) (*model.VerificationCode, error) {
	vCodes, err := s.store.GetVCodesByEmail(email)
	if err != nil {
		return nil, err
	}
	if len(vCodes) == 5 {
		vCodes, err = s.removeExpired(vCodes)
		if err != nil {
			return nil, err
		}

		if len(vCodes) == 5 {
			return nil, NewValidationError("* Can't request more codes. Try again later.", ErrConflict)
		}
	}

	code, err := model.NewVerificationCode(email)
	if err != nil {
		return nil, err
	}

	if err := s.store.CreateVCode(code); err != nil {
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

	exists, err := s.store.EmailExists(user.Email)
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

func (s *userService) Verify(subCode, username, email, password string) error {
	vCodes, err := s.store.GetVCodesByEmail(email)
	if err != nil {
		return err
	}

	verify := false
	for _, code := range vCodes {
		if code.ExpiresAt.After(time.Now()) && subCode == code.Code {
			verify = true
			break
		}
	}

	if !verify {
		return NewValidationError("* Incorrect code, try again.", ErrConflict)
	}

	user := model.NewUser(username, email, password)
	if err := user.HashPassword(); err != nil {
		return err
	}

	if err := s.store.Create(user); err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			return NewValidationError("* Incorrect code, try again", ErrConflict)
		case store.ErrDuplicateUsername:
			return NewValidationErrorsSingle("username", "* Username already taken.", ErrConflict)
		default:
			return err
		}
	}

	return nil
}
