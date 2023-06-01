package models

import (
	"database/sql"
	"fmt"
	"time"
)

type PasswordReset struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB            *sql.DB
	BytesPerToken int
	Duration      time.Duration
}

const (
	DefaultResetDuration = 1 * time.Hour
)

func (svc *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, fmt.Errorf("password Reset Create: Implement me")
}

func (svc *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("password reset consume: Implement me")
}
