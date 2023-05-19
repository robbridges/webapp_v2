package models

import (
	"fmt"
	"github.com/robbridges/webapp_v2/rand"
)

const (
	MinBytesPerToken = 32
)

type tokenManager struct {
	BytesPerToken int
}

func (tm *tokenManager) New() (string, error) {
	bytesPerToken := tm.BytesPerToken
	// Check what Bytes per Token is set 0, if not set or less than the min bytes we override it to the min bytes.
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err := rand.String(bytesPerToken)
	if err != nil {
		return "", fmt.Errorf("create session token: %w", err)
	}
	return token, nil
}
