package models

import (
	"errors"
	"testing"
)

func TestMockSessionService_Create(t *testing.T) {
	mockSession := &Session{UserID: 123, Token: "abc123"}

	// create a new mock session service
	mockSessionService := &MockSessionService{}

	// set the mock Create function to return a session
	mockSessionService.CreateFunc = func(userID int) (*Session, error) {
		return mockSession, nil
	}

	// call the Create method on the mock session service
	session, err := mockSessionService.Create(123)

	// check if the returned session matches the expected session
	if session != mockSession {
		t.Errorf("session does not match the expected session")
	}

	// check if the error is nil
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMockSessionService_User(t *testing.T) {
	// Set up mock user
	expectedUser := &User{ID: 1, Email: "test@test.com"}

	// Set up mock session service
	mockSessionService := &MockSessionService{}

	// Test case 1: UserFunc is not set
	user, err := mockSessionService.User("abc123")
	if user != nil || err == nil {
		t.Errorf("unexpected result: user=%v, err=%v", user, err)
	}

	// Test case 2: UserFunc is set
	mockSessionService.UserFunc = func(token string) (*User, error) {
		if token == "abc123" {
			return expectedUser, nil
		}
		return nil, errors.New("invalid token")
	}
	user, err = mockSessionService.User("abc123")
	if user != expectedUser || err != nil {
		t.Errorf("unexpected result: user=%v, err=%v", user, err)
	}

	// Test case 3: Invalid token
	user, err = mockSessionService.User("invalid_token")
	if user != nil || err == nil {
		t.Errorf("unexpected result: user=%v, err=%v", user, err)
	}
}

func TestDeleteSession(t *testing.T) {
	mockService := &MockSessionService{}

	testCases := []struct {
		name        string
		token       string
		expectedErr string
	}{
		{
			name:        "Token found",
			token:       "found token",
			expectedErr: "",
		},
		{
			name:        "Token not found",
			token:       "not found token",
			expectedErr: "this token was not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := mockService.DeleteSession(tc.token)
			if err == nil && tc.expectedErr != "" {
				t.Errorf("Expected error '%s', but got no error", tc.expectedErr)
			} else if err != nil && err.Error() != tc.expectedErr {
				t.Errorf("Expected error '%s', but got '%s'", tc.expectedErr, err.Error())
			}
		})
	}
}
