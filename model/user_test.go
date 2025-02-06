package model

import (
	"errors"
	"net/http"
	"testing"
)

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

type MockUserStore struct {
	UserNameExistsFn func(string) (bool, error)
}

func (m *MockUserStore) UserNameExists(u string) (bool, error) {
	return m.UserNameExistsFn(u)
}

func (m *MockUserStore) Create(u *User) StatusCoder {
	return nil
}

func (m *MockUserStore) CreateVerificationCode(v *VerificationCode) error {
	return nil
}

func (m *MockUserStore) GetUserByID(id uint) (*User, error) {
	return nil, nil
}

func TestValidUsername(t *testing.T) {
	tests := []struct {
		name            string
		username        string
		mockUserNameFn  func(string) (bool, error)
		expectedStatus  int
		expectedMessage string
	}{
		{"required username", "", func(u string) (bool, error) { return false, nil }, http.StatusUnprocessableEntity, "* Username is required."},
		{"valid username underscore", "ValidUser_1", func(u string) (bool, error) { return false, nil }, http.StatusOK, ""},
		{"valid username dash", "ValidUser-1", func(u string) (bool, error) { return false, nil }, http.StatusOK, ""},
		{"len3", "123", func(u string) (bool, error) { return false, nil }, http.StatusOK, ""},
		{"len2", "12", func(u string) (bool, error) { return false, nil }, http.StatusUnprocessableEntity, "* Username must have at least 3 characters."},
		{"len15", "123456789ABCDEF", func(u string) (bool, error) { return false, nil }, http.StatusOK, ""},
		{"len16", "123456789ABCDEFG", func(u string) (bool, error) { return false, nil }, http.StatusUnprocessableEntity, "* Username must have at most 15 characters."},
		{"username with @", "Invalid@User", func(u string) (bool, error) { return false, nil }, http.StatusUnprocessableEntity, "* English letters, digits, _ and - only."},
		{"username with space", "Invalid User", func(u string) (bool, error) { return false, nil }, http.StatusUnprocessableEntity, "* English letters, digits, _ and - only."},
		{"username with *", "Invalid*User", func(u string) (bool, error) { return false, nil }, http.StatusUnprocessableEntity, "* English letters, digits, _ and - only."},
		{"username already taken", "TakenUser", func(u string) (bool, error) { return true, nil }, http.StatusConflict, "* Username already taken."},
		{"db error", "TakenUser", func(u string) (bool, error) { return false, errors.New("") }, http.StatusInternalServerError, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &MockUserStore{UserNameExistsFn: tt.mockUserNameFn}

			result := ValidUsername(tt.username, us)

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
