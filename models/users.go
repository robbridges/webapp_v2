package models

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	passwordHash := string(hashedBytes)

	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}
	row := us.DB.QueryRow(`
		INSERT INTO USERS (email, password_hash)
		VALUES ($1, $2) RETURNING id;`, email, passwordHash,
	)
	if err := row.Scan(&user.ID); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	//stubbed out, but we'd convert this email and password to a new User to save to the database.
	return &user, nil
}
