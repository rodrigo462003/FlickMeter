package model

import (
	"errors"
	"net/http"
	"testing"
)

type MockUserStore struct {
	UserNameExistsFunc         func(string) (bool, error)
	CreateFunc                 func(*User) StatusCoder
	CreateVerificationCodeFunc func(*VerificationCode) error
	GetUserByIDFunc            func(uint) (*User, error)
}

func (m *MockUserStore) UserNameExists(username string) (bool, error) {
	return m.UserNameExistsFunc(username)
}

func (m *MockUserStore) Create(u *User) StatusCoder {
	return m.CreateFunc(u)
}

func (m *MockUserStore) CreateVerificationCode(v *VerificationCode) error {
	return m.CreateVerificationCodeFunc(v)
}

func (m *MockUserStore) GetUserByEmail(email string) (user *User, err error) {
	return nil, nil
}

func (m *MockUserStore) GetUserByID(id uint) (*User, error) {
	return m.GetUserByIDFunc(id)
}

func (m *MockUserStore) DeleteCode(v *VerificationCode) error {
	return nil
}

type MockEmailSender struct{}

func (m *MockEmailSender) SendMail(to, subject, body string) {
	return
}

func TestNewUser(t *testing.T) {
	tests := []struct {
		name            string
		username        string
		email           string
		password        string
		mockStore       *MockUserStore
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:     "successful user creation",
			username: "validUser",
			email:    "user@example.com",
			password: "password123",
			mockStore: &MockUserStore{
				CreateVerificationCodeFunc: func(vc *VerificationCode) error { return nil },
				UserNameExistsFunc:         func(username string) (bool, error) { return false, nil },
				CreateFunc:                 func(u *User) StatusCoder { return nil },
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "",
		},
		{
			name:     "username already exists",
			username: "existingUser",
			email:    "existing@example.com",
			password: "password123",
			mockStore: &MockUserStore{
				CreateVerificationCodeFunc: func(vc *VerificationCode) error { return nil },
				UserNameExistsFunc:         func(username string) (bool, error) { return true, nil },
				CreateFunc:                 func(u *User) StatusCoder { return nil },
			},
			expectedStatus:  http.StatusConflict,
			expectedMessage: "map[username:* Username already taken.]",
		},
		{
			name:     "validation error (username too short)",
			username: "sh",
			email:    "invalid@example.com",
			password: "123",
			mockStore: &MockUserStore{
				CreateVerificationCodeFunc: func(vc *VerificationCode) error { return nil },
				UserNameExistsFunc:         func(username string) (bool, error) { return false, nil },
				CreateFunc:                 func(u *User) StatusCoder { return nil },
			},
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedMessage: "map[password:* Must contain at least 8 characters. username:* Username must have at least 3 characters.]",
		},
		{
			name:     "internal server error during user creation",
			username: "validUser",
			email:    "user@example.com",
			password: "password123",
			mockStore: &MockUserStore{
				CreateVerificationCodeFunc: func(vc *VerificationCode) error { return nil },
				UserNameExistsFunc:         func(username string) (bool, error) { return false, nil },
				CreateFunc: func(u *User) StatusCoder {
					return NewStatusError(http.StatusInternalServerError, "Something unexpected has happened, please try again.")
				},
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Something unexpected has happened, please try again.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &MockEmailSender{}
			result := NewUser(tt.username, tt.email, tt.password, tt.mockStore, es)

			if result != nil {
				if result.StatusCode() != tt.expectedStatus {
					t.Errorf("Expected status code %d, got %d", tt.expectedStatus, result.StatusCode())
				}
				if result.Error() != tt.expectedMessage {
					t.Errorf("Expected message %q, got %q", tt.expectedMessage, result.Error())
				}
			} else {
				if tt.expectedStatus != http.StatusOK {
					t.Errorf("Expected status code %d, got %d", http.StatusOK, tt.expectedStatus)
				}
				if tt.expectedMessage != "" {
					t.Errorf("Expected empty message, got %q", tt.expectedMessage)
				}
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name        string
		user        *User
		userStore   *MockUserStore
		expectedErr StatusCoder
	}{
		{"Valid user", &User{Username: "validuser", Email: "valid@example.com", Password: "ValidPass123"}, &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, nil},
		{"Invalid username", &User{Username: "invalid_user@!", Email: "valid@example.com", Password: "ValidPass123"}, &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, NewStatusErrors(http.StatusUnprocessableEntity, map[string]string{"username": "* English letters, digits, _ and - only."})},
		{"Username already exists", &User{Username: "existinguser", Email: "valid@example.com", Password: "ValidPass123"}, &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return true, nil }}, NewStatusErrors(http.StatusConflict, map[string]string{"username": "* Username already taken."})},
		{"Invalid email", &User{Username: "validuser", Email: "invalidemail", Password: "ValidPass123"}, &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, NewStatusErrors(http.StatusUnprocessableEntity, map[string]string{"email": "* This is not a valid email address."})},
		{"Invalid password", &User{Username: "validuser", Email: "valid@example.com", Password: "short"}, &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, NewStatusErrors(http.StatusUnprocessableEntity, map[string]string{"password": "* Must contain at least 8 characters."})},
		{"All invalid", &User{Username: "val432**", Email: "validexample.com", Password: "sho"}, &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, NewStatusErrors(http.StatusUnprocessableEntity, map[string]string{"password": "* Must contain at least 8 characters.", "username": "* English letters, digits, _ and - only.", "email": "* This is not a valid email address."})},
		{"One valid", &User{Username: "val432**", Email: "valide@xample.com", Password: "sho"}, &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, NewStatusErrors(http.StatusUnprocessableEntity, map[string]string{"password": "* Must contain at least 8 characters.", "username": "* English letters, digits, _ and - only."})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.isValid(tt.userStore)

			if tt.expectedErr != nil {
				var got StatusErrors
				ok := errors.As(err, &got)
				if !ok {
					t.Error("isValid() errors.As(err, StatusErrors) failed")
					return
				}
				var expected StatusErrors
				ok = errors.As(tt.expectedErr, &expected)
				if !ok {
					t.Error("isValid() errors.As(tt.expectedErr, StatusErrors) failed")
					return
				}

				if err == nil || got.StatusCode() != tt.expectedErr.StatusCode() {
					t.Errorf("isValid() = %v, want %v", got, tt.expectedErr)
				}
				if err != nil && got.errorMap != nil {
					for field, expectedMessage := range expected.Map() {
						if gotMessage, exists := got.errorMap[field]; !exists || gotMessage != expectedMessage {
							t.Errorf("isValid() error for %s = %v, want %v", field, gotMessage, expectedMessage)
						}
					}
				}
			} else if err != nil {
				t.Errorf("isValid() = %v, want nil", err)
			}
		})
	}
}

func TestGetPriorityStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"No matching codes", []int{http.StatusOK, http.StatusCreated}, http.StatusInternalServerError},
		{"One code", []int{http.StatusConflict}, http.StatusConflict},
		{"conflict", []int{http.StatusOK, http.StatusConflict}, http.StatusConflict},
		{"internal over unprocess", []int{http.StatusUnprocessableEntity, http.StatusInternalServerError}, http.StatusInternalServerError},
		{"All codes included", []int{http.StatusInternalServerError, http.StatusConflict, http.StatusUnprocessableEntity}, http.StatusInternalServerError},
		{"Conflict over unprocess", []int{http.StatusConflict, http.StatusUnprocessableEntity}, http.StatusConflict},
		{"internal over conflict", []int{http.StatusInternalServerError, http.StatusConflict}, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := getPriorityStatusCode(tt.input)
			if got != tt.expected {
				t.Errorf("getPriorityStatusCode() = %v, want %v", got, tt.expected)
			}

		})
	}
}

func TestValidUsername(t *testing.T) {
	tests := []struct {
		name            string
		username        string
		userStore       *MockUserStore
		expectedStatus  int
		expectedMessage string
	}{
		{"required username", "", &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, http.StatusUnprocessableEntity, "* Username is required."},
		{"valid username underscore", "ValidUser_1", &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, http.StatusOK, ""},
		{"valid username dash", "ValidUser-1", &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, http.StatusOK, ""},
		{"len3", "123", &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, http.StatusOK, ""},
		{"len2", "12", &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, http.StatusUnprocessableEntity, "* Username must have at least 3 characters."},
		{"len15", "123456789ABCDEF", &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, http.StatusOK, ""},
		{"len16", "123456789ABCDEFG", &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, http.StatusUnprocessableEntity, "* Username must have at most 15 characters."},
		{"username with @", "Invalid@User", &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, http.StatusUnprocessableEntity, "* English letters, digits, _ and - only."},
		{"username with space", "Invalid User", &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, http.StatusUnprocessableEntity, "* English letters, digits, _ and - only."},
		{"username with *", "Invalid*User", &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, nil }}, http.StatusUnprocessableEntity, "* English letters, digits, _ and - only."},
		{"username already taken", "TakenUser", &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return true, nil }}, http.StatusConflict, "* Username already taken."},
		{"db error", "TakenUser", &MockUserStore{UserNameExistsFunc: func(username string) (bool, error) { return false, errors.New("") }}, http.StatusInternalServerError, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidUsername(tt.username, tt.userStore)

			if result != nil {
				if result.StatusCode() != tt.expectedStatus {
					t.Errorf("Expected status code %d, got %d", tt.expectedStatus, result.StatusCode())
				}
				if result.Error() != tt.expectedMessage {
					t.Errorf("Expected message %q, got %q", tt.expectedMessage, result.Error())
				}
			} else {
				if tt.expectedStatus != http.StatusOK {
					t.Errorf("Expected status code %d, got %d", http.StatusOK, tt.expectedStatus)
				}
				if tt.expectedMessage != "" {
					t.Errorf("Expected empty message, got %q", tt.expectedMessage)
				}
			}
		})
	}
}

func TestValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		wantCode int
		wantMsg  string
	}{
		{"EmptyEmail", "", http.StatusUnprocessableEntity, "* Email address is required."},
		{"InvalidEmail_NoAt", "invalidemail.com", http.StatusUnprocessableEntity, "* This is not a valid email address."},
		{"InvalidEmail_NoDomain", "invalid@", http.StatusUnprocessableEntity, "* This is not a valid email address."},
		{"InvalidEmail_Spaces", "invalid email@example.com", http.StatusUnprocessableEntity, "* This is not a valid email address."},
		{"ValidEmail", "valid@example.com", 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidEmail(tt.email)
			if got != nil {
				if got.StatusCode() != tt.wantCode || got.Error() != tt.wantMsg {
					t.Errorf("%s: ValidEmail() = %v, want code %d and message %q", tt.name, got, tt.wantCode, tt.wantMsg)
				}
			} else if tt.wantCode != 0 {
				t.Errorf("%s: ValidEmail() = nil, want code %d and message %q", tt.name, tt.wantCode, tt.wantMsg)
			}
		})
	}
}

func TestValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantCode int
		wantMsg  string
	}{
		{"empty", "", http.StatusUnprocessableEntity, "* Password is required."},
		{"len7", "1234567", http.StatusUnprocessableEntity, "* Must contain at least 8 characters."},
		{"len8", "12345678", 0, ""},
		{"len9", "123456789", 0, ""},
		{"len127", "thispasswordiswaytoolongandexceeds128character23432432423423333333333333333333333333333333333333333334322434242343243242342s2..", 0, ""},
		{"len128", "thispasswordiswaytoolongandexceeds128character23432432423423333333333333333333333333333333333333333334322434242343243242342s2...", 0, ""},
		{"len129", "thispasswordiswaytoolongandexceeds128character23432432423423333333333333333333333333333333333333333334322434242343243242342s2....", http.StatusUnprocessableEntity, "* Must contain at most 128 characters."},
		{"invalid_utf8", string([]byte{0xff, 0xfe, 0xfd}), http.StatusUnprocessableEntity, "* Invalid character(s) detected, try again."},
		{"multi_codepoint_len7", "🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿", http.StatusUnprocessableEntity, "* Must contain at least 8 characters."},
		{"multi_codepoint_len8", "🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿", 0, ""},
		{"multi_codepoint_len9", "🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿🏴󠁧󠁢󠁳󠁣󠁴󠁿", 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidPassword(tt.password)

			if got != nil {
				if got.StatusCode() != tt.wantCode || got.Error() != tt.wantMsg {
					t.Errorf("ValidPassword(%q) = %v, want code %d and message %q", tt.password, got, tt.wantCode, tt.wantMsg)
				}
			} else if tt.wantCode != 0 {
				t.Errorf("ValidPassword(%q) = nil, want code %d and message %q", tt.password, tt.wantCode, tt.wantMsg)
			}
		})
	}
}
