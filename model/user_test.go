package model

import (
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
		{"len7", "1234567", http.StatusUnprocessableEntity, "* Password must contain atleast 8 characters."},
		{"len8", "12345678", 0, ""},
		{"len9", "123456789", 0, ""},
		{"len127", "thispasswordiswaytoolongandexceeds128character23432432423423333333333333333333333333333333333333333334322434242343243242342s2..", 0, ""},
		{"len128", "thispasswordiswaytoolongandexceeds128character23432432423423333333333333333333333333333333333333333334322434242343243242342s2...", 0, ""},
		{"len129", "thispasswordiswaytoolongandexceeds128character23432432423423333333333333333333333333333333333333333334322434242343243242342s2....", http.StatusUnprocessableEntity, "* Password must contain at most 128 characters."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidPassword(tt.password)

			if got != nil {
				if got.StatusCode() != tt.wantCode || got.Error() != tt.wantMsg {
					t.Errorf("ValidPassword() = %v, want code %d and message %q", got, tt.wantCode, tt.wantMsg)
				}
			} else if tt.wantCode != 0 {
				t.Errorf("ValidPassword() = nil, want code %d and message %q", tt.wantCode, tt.wantMsg)
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
