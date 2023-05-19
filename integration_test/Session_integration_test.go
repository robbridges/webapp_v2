package integration_test

import (
	"github.com/robbridges/webapp_v2/models"
	"testing"
	"time"
)

func TestCreateSession(t *testing.T) {

	// Set up a test database
	db, err := setup(t)

	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer deferDBClose(db, &err)
	defer teardown()

	sessionService := &models.SessionService{DB: db}

	err = waitForPing(db, 10*time.Second)
	if err != nil {
		t.Errorf("Database timeout: %v", err)
	}

	userID := 1
	session, err := sessionService.Create(userID)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	if session.ID == 0 {
		t.Error("session ID is not set")
	}
	if session.UserID != userID {
		t.Errorf("expected user ID %d, got %d", userID, session.UserID)
	}
	if len(session.Token) == 0 {
		t.Error("session token is not set")
	}

	var storedSession models.Session
	err = db.QueryRow("SELECT * FROM sessions WHERE id = $1", session.ID).Scan(&storedSession.ID, &storedSession.UserID, &storedSession.TokenHash)
	if err != nil {
		t.Fatalf("failed to query stored session: %v", err)
	}
	if storedSession.ID != session.ID {
		t.Errorf("expected stored session ID %d, got %d", session.ID, storedSession.ID)
	}
	if storedSession.UserID != session.UserID {
		t.Errorf("expected stored user ID %d, got %d", session.UserID, storedSession.UserID)
	}
	if storedSession.TokenHash != models.Hash(session.Token) {
		t.Error("stored session token hash is incorrect")
	}

	session2, err := sessionService.Create(userID)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	if session2.Token == session.Token {
		t.Error("new session has the same token as the previous session")
	}

	err = db.QueryRow("SELECT * FROM sessions WHERE id = $1", session2.ID).Scan(&storedSession.ID, &storedSession.UserID, &storedSession.TokenHash)
	if err != nil {
		t.Fatalf("failed to query stored session: %v", err)
	}
	if storedSession.ID != session2.ID {
		t.Errorf("expected stored session ID %d, got %d", session2.ID, storedSession.ID)
	}
	if storedSession.UserID != session2.UserID {
		t.Errorf("expected stored user ID %d, got %d", session2.UserID, storedSession.UserID)
	}
	if storedSession.TokenHash != models.Hash(session2.Token) {
		t.Error("stored session token hash is incorrect")
	}
}

func TestGetUserByToken(t *testing.T) {
	db, err := setup(t)

	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	defer deferDBClose(db, &err)
	defer teardown()

	sessionService := &models.SessionService{DB: db}

	err = waitForPing(db, 10*time.Second)
	if err != nil {
		t.Errorf("Database timeout: %v", err)
	}

	userID := 1
	email := "test@example.com"
	passwordHash := "testpasswordhash"
	_, err = db.Exec("INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3)", userID, email, passwordHash)
	if err != nil {
		t.Fatalf("failed to set up test user: %v", err)
	}

	// Set up a test session
	token := "testtoken"
	tokenHash := models.Hash(token)
	_, err = db.Exec("INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2)", userID, tokenHash)
	if err != nil {
		t.Fatalf("failed to set up test session: %v", err)
	}

	user, err := sessionService.User(token)
	if err != nil {
		t.Fatalf("failed to get user by token: %v", err)
	}

	if user.Email != email {
		t.Errorf("expected user email %s, got %s", email, user.Email)
	}
	if user.PasswordHash != passwordHash {
		t.Error("user password hash is incorrect")
	}

	invalidToken := "invalidtoken"
	_, err = sessionService.User(invalidToken)
	if err == nil {
		t.Error("expected an error when getting user by invalid token")
	}
}

func TestDeleteSession(t *testing.T) {
	db, err := setup(t)

	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	defer deferDBClose(db, &err)
	defer teardown()

	sessionService := &models.SessionService{DB: db}

	err = waitForPing(db, 10*time.Second)
	if err != nil {
		t.Errorf("Database timeout: %v", err)
	}

	token := "testtoken"
	tokenHash := models.Hash(token)
	_, err = db.Exec("INSERT INTO sessions (token_hash) VALUES ($1)", tokenHash)
	if err != nil {
		t.Fatalf("failed to set up test session: %v", err)
	}

	err = sessionService.DeleteSession(token)
	if err != nil {
		t.Fatalf("failed to delete session: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sessions WHERE token_hash = $1", tokenHash).Scan(&count)
	if err != nil {
		t.Fatalf("failed to query session count: %v", err)
	}

	if count != 0 {
		t.Error("session was not deleted from the database")
	}
}
