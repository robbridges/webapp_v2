package integration_test

import (
	"github.com/robbridges/webapp_v2/models"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestUserService_Create(t *testing.T) {
	// Get a database connection
	db, err := setup(t)

	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer deferDBClose(db, &err)
	defer teardown(t)

	us := &models.UserService{DB: db}

	email := "test@test.com"
	password := "password"

	err = waitForPing(db, 10*time.Second)
	if err != nil {
		t.Errorf("Database never responded")
	}

	u, err := us.Create(email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if u.Email != email {
		t.Errorf("Expected email to be %s, got %s", email, u.Email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		t.Errorf("Expected password hashes to match, but they did not: %v", err)
	}

	var dbUser models.User
	err = db.QueryRow("SELECT id, email, password_hash FROM USERS WHERE id = $1", u.ID).Scan(&dbUser.ID, &dbUser.Email, &dbUser.PasswordHash)
	if err != nil {
		t.Fatalf("Failed to fetch user from DB: %v", err)
	}

	if dbUser.Email != email {
		t.Errorf("Expected DB user email to be %s, got %s", email, dbUser.Email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(password))
	if err != nil {
		t.Errorf("Expected DB password hashes to match, but they did not: %v", err)
	}
}

func TestUserService_Authenticate(t *testing.T) {
	// Get a database connection
	db, err := setup(t)

	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	defer deferDBClose(db, &err)
	defer teardown(t)

	us := &models.UserService{DB: db}

	err = waitForPing(db, 10*time.Second)
	if err != nil {
		t.Errorf("Database never responded")
	}

	email := "testauth@test.com"
	password := "authpassword"

	// Before authenticating, let's create a user
	_, err = us.Create(email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	t.Run("Happy Path", func(t *testing.T) {
		// Authenticate the created user
		u, err := us.Authenticate(email, password)
		if err != nil {
			t.Fatalf("Failed to authenticate user: %v", err)
		}

		if u.Email != email {
			t.Errorf("Expected email to be %s, got %s", email, u.Email)
		}

		// Verify password hash
		err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
		if err != nil {
			t.Errorf("Expected password hashes to match, but they did not: %v", err)
		}
	})

	t.Run("Sad Path", func(t *testing.T) {
		// Try to authenticate with wrong password
		_, err = us.Authenticate(email, "wrongpassword")
		if err == nil {
			t.Errorf("Expected error when authenticating with wrong password, got nil")
		}
	})
}
