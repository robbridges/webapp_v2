package models

import (
	"errors"
	"strings"
	"testing"
)

func TestMockUserService_Create(t *testing.T) {
	t.Run("Mock service create - happy", func(t *testing.T) {
		email := "Test@Test.com"
		password := "very secure"
		mus := MockUserService{}
		user, err := mus.Create(email, password)
		if err != nil {
			t.Errorf("Error was returned unexpectedly, %v", err)
		}
		if user.Email != strings.ToLower(email) {
			t.Errorf("email not correct set, wanted %s, got %s", email, user.Email)
		}

		if user.PasswordHash == password {
			t.Errorf("The password was never hashed")
		}
	})
	t.Run("Mock service create - sad", func(t *testing.T) {
		email := "Test@Test.com"
		password := "this_is_a_very_long_password_string_that_exceeds_the_bcrypt_limit_of_72_characters"
		mus := MockUserService{}
		_, err := mus.Create(email, password)
		if err == nil {
			t.Error("no error returned, but was expected")
		}
		got := err.Error()
		want := "failed to hash password: bcrypt: password length exceeds 72 bytes"
		if got != want {
			t.Errorf("We did not get the expected error back, got %s, want %s", got, want)
		}

	})
}

func TestMockUserService_Authenticate(t *testing.T) {
	email := "test@test.com"
	password := "password"

	// Set up mock user service
	mockUserService := &MockUserService{}

	// Test case 1: AuthenticateFunc is not set
	user, err := mockUserService.Authenticate(email, password)
	if user != nil || err == nil {
		t.Errorf("unexpected result: user=%v, err=%v", user, err)
	}

	// Test case 2: AuthenticateFunc is set
	expectedUser := &User{ID: 1, Email: email}
	mockUserService.AuthenticateFunc = func(email, password string) (*User, error) {
		if password == "password" {
			return expectedUser, nil
		}
		return nil, errors.New("invalid credentials")
	}
	user, err = mockUserService.Authenticate(email, password)
	if user != expectedUser || err != nil {
		t.Errorf("unexpected result: user=%v, err=%v", user, err)
	}

	// Test case 3: Invalid credentials
	user, err = mockUserService.Authenticate(email, "wrongpassword")
	if user != nil || err == nil {
		t.Errorf("unexpected result: user=%v, err=%v", user, err)
	}
}
