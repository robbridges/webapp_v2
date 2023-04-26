package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/robbridges/webapp_v2/rand"
)

const (
	// MinBytesPerToken The minimum number of bytes to be used for each Session token.
	MinBytesPerToken = 32
)

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
	// BytesPerToken is used to determine how many bytes to use when generating each session token If this value is
	// under MinBytes Per Token it will be ignored and MinBytesPerToken will be used. MinBytsPerToken will also be used
	// if this value is not set. Just a bit future proofing
	BytesPerToken int
}

// Create will create a new session for the user provided the session token is the returned string to be stored
// in our Postgres user table
func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	// Check what Bytes per Token is set 0, if not set or less than the min bytes we over ride it to the min bytes.
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create session token: %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
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
	//Todo implement this as well
	return nil, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
