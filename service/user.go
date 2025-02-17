package service

import (
	"time"

	"github.com/rodrigo462003/FlickMeter/email"
	"github.com/rodrigo462003/FlickMeter/model"
	"github.com/rodrigo462003/FlickMeter/store"
)

type UserService interface {
	ValidatePassword(password string) error
	ValidateEmail(email string) error
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

func (s *userService) ValidatePassword(password string) error {
	if err := model.NewPassword(password); err != nil {
		return NewValidationError(err.Error(), ErrUnprocessable)
	}

	return nil
}

func (s *userService) ValidateEmail(email string) error {
	if err := model.NewEmail(email); err != nil {
		return NewValidationError(err.Error(), ErrUnprocessable)
	}

	return nil
}

func (s *userService) ValidateUsername(username string) error {
	if err := model.NewUsername(username); err != nil {
		return NewValidationError(err.Error(), ErrUnprocessable)
	}

	isDupe, err := s.store.UsernameExists(username)
	if err != nil {
		return err
	}
	if isDupe {
		return NewValidationError(err.Error(), ErrConflict)
	}

	return nil
}

func (s *userService) validateUser(user *model.User) error {
	errMap := make(map[string]string, 3)
	if err := s.ValidateUsername(user.Username); err != nil {
		errMap["username"] = err.Error()
	}
	if err := s.ValidateEmail(user.Email); err != nil {
		errMap["email"] = err.Error()
	}
	if err := s.ValidatePassword(user.Password); err != nil {
		errMap["password"] = err.Error()
	}

	if len(errMap) > 0 {
		return NewValidationErrors()
	}

	return nil
}

func (s *userService) CreateUser(username, email, password string) error {
	const emailSubject = "FlickMeter registration"

	user := model.NewUser(username, email, password)
	if err := s.validateUser(user); err != nil {
		if vErr, ok := err.(service.ValidationError); ok {
			return
		}
	}

	tmpU := *user
	created, err := s.store.FirstOrCreate(user)
	if err != nil {
		return InternalServerError()
	}
	if !created {
		if user.Verified {
			if user.Username == tmpU.Username && user.Email != tmpU.Email {
				return NewStatusErrors(http.StatusConflict, map[string]string{"username": "Username already taken."})
			}
			if user.Email == tmpU.Email {
				s.sender.SendMail(User.Email, emailSubject, "You already have a FlickMeter account, try signing in.")
				return nil
			}
			slog.Error(err.Error())
			return InternalServerError()
		}

		if user.Username == tmpU.Username && User.Email != tmpU.Email {
			return NewStatusErrors(http.StatusConflict, map[string]string{"username": "Username already taken."})
		}
		if user.Email == tmpU.Email {
			user.Username, user.Password = user.Username, user.Password
			if err := s.store.Save(user); err != nil {
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
	es.SendMail(User.Email, emailSubject, emailBody)

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
