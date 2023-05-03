package models

import (
	"strings"
	"testing"
)

func TestMockUserService_Create(t *testing.T) {
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
}

func TestMockUserService_Authenticate(t *testing.T) {
	t.Run("Mock Authentication happy path", func(t *testing.T) {
		email := "test@admin.com"
		password := "secure"
		mus := MockUserService{}
		user, err := mus.Authenticate(email, password)
		if err != nil {
			t.Errorf("User password %s, should have returned user it returned an error", password)
		}

		if user.Email != email {
			t.Errorf("User did not get assigned the correct details")
		}
	})
	t.Run("Mock Authentication, sad path", func(t *testing.T) {
		mockErrorMessage := "invalid email or password"
		email := "test@admin.com"
		password := "fake"
		mus := MockUserService{}
		_, err := mus.Authenticate(email, password)
		if err == nil {
			t.Errorf("User password %s, should have returned error, it did not", password)
		}
		got := err.Error()
		want := mockErrorMessage
		if got != want {
			t.Errorf("The wrong error was returned we got %s, but wanted %s", got, want)
		}
	})
}
