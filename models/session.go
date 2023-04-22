package models

import "database/sql"

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

// Create will create a new session for the user provided the session token is the returned string to be stored
// in our Postgres user table
func (ss *SessionService) Create(userId int) (*Session, error) {
	//Todo create the session token
	//Todo Implement this method

	return nil, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	//Todo implement this as well
	return nil, nil
}
