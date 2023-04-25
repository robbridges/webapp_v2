package models

import (
	"database/sql"
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
	//TODO hash session token
	if err != nil {
		return nil, fmt.Errorf("create session token: %w", err)
	}
	session := Session{
		UserID: userID,
		Token:  token,
		//TODO set token hash
	}

	//Todo store session in database
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	//Todo implement this as well
	return nil, nil
}
