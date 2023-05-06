package models

import (
	"errors"
	"testing"
)

func TestUserService_Create(t *testing.T) {
	mockUserService := &MockUserService{}

	t.Run("happy path", func(t *testing.T) {
		mockUserService.CreateFunc = func(email string, password string) (*User, error) {
			user := &User{
				ID:           1,
				Email:        email,
				PasswordHash: "hashedPassword",
			}
			return user, nil
		}

		user, err := mockUserService.Create("test@test.com", "password")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if user.ID != 1 {
			t.Errorf("unexpected user ID: got %v, want %v", user.ID, 1)
		}
		if user.Email != "test@test.com" {
			t.Errorf("unexpected email: got %v, want %v", user.Email, "test@test.com")
		}
		if user.PasswordHash != "hashedPassword" {
			t.Errorf("unexpected password hash: got %v, want %v", user.PasswordHash, "hashedPassword")
		}
	})

	t.Run("error path", func(t *testing.T) {
		mockUserService.CreateFunc = func(email string, password string) (*User, error) {
			return nil, errors.New("database error")
		}

		_, err := mockUserService.Create("test@test.com", "password")
		expectedErr := "database error"
		if err == nil || err.Error() != expectedErr {
			t.Errorf("unexpected error: got %v, want %v", err, expectedErr)
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
