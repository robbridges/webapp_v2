package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/robbridges/webapp_v2/rand"
	"strings"
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
	DB *sql.DB
	//Same as the session service, if this value is not set, or is less than the min constant in the Token package it
	// will automatically be set to 32
	BytesPerToken int
	Duration      time.Duration
}

const (
	DefaultResetDuration = 1 * time.Hour
)

func (svc *PasswordResetService) Create(email string) (*PasswordReset, error) {
	email = strings.ToLower(email)

	var userID int
	row := svc.DB.QueryRow(`
	SELECT id FROM users where email = $1;`, email)
	err := row.Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("create %v", err)
	}
	// build password reset token
	bytesPerToken := svc.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create %v", err)
	}
	duration := svc.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}

	pwReset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: svc.hash(token),
		ExpiresAt: time.Now().Add(duration),
	}

	row = svc.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id;`, pwReset.UserID, pwReset.TokenHash, pwReset.ExpiresAt)
	err = row.Scan(&pwReset.ID)

	if err != nil {
		return nil, fmt.Errorf("create %w", err)
	}

	return &pwReset, nil
}

func (svc *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("password reset consume: Implement me")
}

func (svc *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
