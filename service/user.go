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
	CreateUser(username, email, password string) error
	VerifyUser(code, email string) error
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
	vErrs := make(map[string]ValidationError)
	if err := s.ValidateUsername(user.Username); err != nil {
		if vErr, ok := err.(ValidationError); ok {
			vErrs["username"] = vErr
		} else {
			return err
		}
	}
	if err := s.ValidateEmail(user.Email); err != nil {
		vErrs["email"] = err
	}
	if err := s.ValidatePassword(user.Password); err != nil {
		vErrs["password"] = err
	}

	if len(vErrs) > 0 {
		return NewValidationErrors(vErrs)
	}

	return nil
}

func (s *userService) removeExpired(user *model.User) error {
	user.RemoveExpiredCodes()
	if err := s.store.Save(user); err != nil {
		return err
	}

	return nil
}

func (s *userService) createVerificationCode(user *model.User) error {
	if len(user.VerificationCodes) == 5 {
		err := s.removeExpired(user)
		if err != nil {
			return err
		}
		if len(user.VerificationCodes) == 5 {
			return NewValidationError("* Can't request more codes. Try again later.", ErrConflict)
		}
	}

	if err := user.NewVerificationCode(); err != nil {
		return err
	}

	if err := s.store.Save(user); err != nil {
		return err
	}

	return nil
}

func (s *userService) handleExistingUser(existingUser *model.User, newUser *model.User) error {
	if existingUser.Verified {
		if existingUser.Username == newUser.Username && existingUser.Email != newUser.Email {
			return NewValidationErrors(map[string]ValidationError{
				"username": NewValidationError("* Username already taken.", ErrConflict),
			})
		}
		if existingUser.Email == newUser.Email {
			return NewValidationError("*User already verified.", errAlreadyVerified)
		}
		panic("Unexpected state: verified user has a username and email conflict")
	}

	if existingUser.Username == newUser.Username && existingUser.Email != newUser.Email {
		return NewValidationErrors(map[string]ValidationError{
			"username": NewValidationError("* Username already taken.", ErrConflict),
		})
	}
	if existingUser.Email == newUser.Email {
		return nil
	}

	panic("Unexpected state: unverified user has a username and email conflict")
}

func (s *userService) CreateUser(username, email, password string) error {
	const (
		emailSubject = "FlickMeter registration"
		codeBodyF    = "Please enter the following code to complete your Signup.\r\n%s"
		verifiedBody = "You already have a FlickMeter account, try signing in."
	)

	user := model.NewUser(username, email, password)
	if err := s.validateUser(user); err != nil {
		return err
	}

	newUser := *user
	created, err := s.store.FirstOrCreate(user)
	if err != nil {
		return err
	}
	if !created {
		if err := s.handleExistingUser(user, &newUser); err != nil {
			if vErr, ok := err.(*validationError); ok && vErr.Is(errAlreadyVerified) {
				s.sender.SendMail(newUser.Email, emailSubject, verifiedBody)
				return nil
			}
			return err
		}
	}

	if err := s.createVerificationCode(user); err != nil {
		return err
	}
	code := user.VerificationCodes[len(user.VerificationCodes)-1].Code

	s.sender.SendMail(user.Email, emailSubject, fmt.Sprintf(codeBodyF, code))
	return nil
}

func (s *userService) VerifyUser(code string, email string) error {
	user, err := s.store.GetByEmail(email)
	if err != nil {
		return err
	}
	if user.Verified {
		return nil
	}

	now := time.Now()
	for _, vCode := range user.VerificationCodes {
		if vCode.ExpiresAt.After(now) && code == vCode.Code {
			user.Verified = true
			return s.store.Save(user)
		}
	}

	return NewValidationError("* Incorrect code.", ErrConflict)
}
