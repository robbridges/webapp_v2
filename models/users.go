package models

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserServiceInterface interface {
	Create(email, password string) (*User, error)
	Authenticate(email, password string) (*User, error)
}

type MockUserService struct {
	AuthenticateFunc func(email, password string) (*User, error)
	CreateFunc       func(email string, password string) (*User, error)
}

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
		return nil, fmt.Errorf("failed to Hash password: %w", err)
	}
	passwordHash := string(hashedBytes)

	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	if err := us.InsertUser(&user); err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &user, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}

	row := us.DB.QueryRow(
		`SELECT id, password_hash
		FROM users WHERE email=$1`, email,
	)

	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		fmt.Errorf("compare() error: %v", err)
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	return &user, nil
}

func (us *UserService) InsertUser(user *User) error {
	row := us.DB.QueryRow(`
		INSERT INTO USERS (email, password_hash)
		VALUES ($1, $2) RETURNING id;`, user.Email, user.PasswordHash,
	)
	if err := row.Scan(&user.ID); err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (mus *MockUserService) Create(email string, password string) (*User, error) {
	if mus.CreateFunc != nil {
		return mus.CreateFunc(email, password)
	}
	email = strings.ToLower(email)

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to Hash password: %w", err)
	}
	passwordHash := string(hashedBytes)

	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	return &user, nil
}

func (mus *MockUserService) Authenticate(email, password string) (*User, error) {
	if mus.AuthenticateFunc != nil {
		return mus.AuthenticateFunc(email, password)
	}
	return nil, errors.New("AuthenticateFunc is not set")
}
