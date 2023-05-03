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
