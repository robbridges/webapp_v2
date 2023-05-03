package models

import (
	"errors"
	"reflect"
	"testing"
)

func TestCreateSession_Success(t *testing.T) {
	mockSessionService := MockSessionService{}

	userID := 123
	session, err := mockSessionService.Create(userID)

	if session == nil {
		t.Errorf("Expected session object, got nil")
	}

	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	if session.UserID != userID {
		t.Errorf("Expected session user ID to be %d, got %d", userID, session.UserID)
	}

	if session.Token == "" {
		t.Errorf("Expected session token to be non-empty, got empty string")
	}

	if session.TokenHash == "" {
		t.Errorf("Expected session token hash to be non-empty, got empty string")
	}
}

func TestMockSessionService_User(t *testing.T) {
	expectedUser := &User{
		ID:           1,
		Email:        "found@user.com",
		PasswordHash: hash("valid_token"),
	}

	// While I'm not certain this will ever get big enough to need table tests, it's easier to do it now then have to
	// later
	testCases := []struct {
		name             string
		token            string
		expectedUser     *User
		expectedError    error
		expectedErrorMsg string
	}{
		{
			name:             "User exists",
			token:            "valid_token",
			expectedUser:     expectedUser,
			expectedError:    nil,
			expectedErrorMsg: "",
		},
		{
			name:             "User does not exist",
			token:            "invalid_token",
			expectedUser:     nil,
			expectedError:    errors.New("no user found by token"),
			expectedErrorMsg: "no user found by token",
		},
	}

	// Loop over each test case and run the test
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the User method of the mock service with the test token
			mockService := MockSessionService{}
			user, err := mockService.User(tc.token)

			// Assert that the returned user object and error match the expected values
			if !reflect.DeepEqual(user, tc.expectedUser) {
				t.Errorf("Expected user to be %v, but got %v", tc.expectedUser, user)
			}
			if err != nil && err.Error() != tc.expectedErrorMsg {
				t.Errorf("Expected error message to be '%s', but got '%s'", tc.expectedErrorMsg, err.Error())
			}
		})
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
