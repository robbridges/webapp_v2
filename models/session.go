package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
)

type SessionServiceInterface interface {
	Create(userID int) (*Session, error)
	User(token string) (*User, error)
	DeleteSession(token string) error
}

type Session struct {
	ID     int
	UserID int
	//Token is only set when creating a new session, we only store the has into the db, so if you're looking up a session
	// This will be unavailable
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

type MockSessionService struct {
	CreateFunc        func(userID int) (*Session, error)
	UserFunc          func(token string) (*User, error)
	DeleteSessionFunc func(token string) error
}

// Create will create a new session for the user provided the session token is the returned string to be stored
// in our Postgres user table
func (ss *SessionService) Create(userID int) (*Session, error) {
	tokenManager := tokenManager{
		BytesPerToken: 32,
	}
	token, err := tokenManager.New()
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: Hash(token),
	}
	row := ss.DB.QueryRow(`
		UPDATE sessions
		SET token_hash = $2
		WHERE user_id = $1
    	RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		// If no session exists, we will get ErrNoRows. That means we need to
		// create a session object for that user.
		row = ss.DB.QueryRow(`
			INSERT INTO sessions (user_id, token_hash)
			VALUES ($1, $2)
			RETURNING id;`, session.UserID, session.TokenHash)
		// The error will be overwritten with either a new error, or nil
		err = row.Scan(&session.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("create %w", err)
	}

	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := Hash(token)
	var user User
	row := ss.DB.QueryRow(`
	SELECT u.email, u.password_hash 
	FROM sessions s
	INNER JOIN users u ON s.user_id = u.id
	WHERE s.token_hash = $1;
	`, tokenHash)
	if err := row.Scan(&user.Email, &user.PasswordHash); err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	return &user, nil
}

func (ss *SessionService) DeleteSession(token string) error {
	tokenHash := Hash(token)
	_, err := ss.DB.Exec(`
	DELETE FROM sessions 
	WHERE token_hash = $1
	`, tokenHash)
	if err != nil {
		return fmt.Errorf("Signout: %w", err)
	}
	return nil
}

func Hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (mss *MockSessionService) Create(userID int) (*Session, error) {
	if mss.CreateFunc != nil {
		return mss.CreateFunc(userID)
	}
	return nil, errors.New("not implemented")
}

func (mss *MockSessionService) User(token string) (*User, error) {
	if mss.UserFunc != nil {
		return mss.UserFunc(token)
	}
	return nil, errors.New("not implemented")
}

func (mss *MockSessionService) DeleteSession(token string) error {
	if mss.DeleteSessionFunc != nil {
		return mss.DeleteSessionFunc(token)
	}
	return errors.New("not implemented")
}
