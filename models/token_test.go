package models

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestTokenManager_New(t *testing.T) {
	t.Run("Proper amount of bytes", func(t *testing.T) {
		tm := setTokenManager(32)

		decodedStr, err := decodeStr(tm)

		if err != nil {
			t.Errorf("Fail: error creating Token: %s", err.Error())
		}

		want := MinBytesPerToken

		got := len(decodedStr)

		if got != want {
			t.Errorf("Got %d: wanted %d", got, want)
		}
	})

	t.Run("Improper amount of bytes", func(t *testing.T) {
		tm := setTokenManager(6)

		decodedStr, err := decodeStr(tm)

		if err != nil {
			t.Errorf("error decoding str %s", err.Error())
		}

		want := MinBytesPerToken

		got := len(decodedStr)

		if got != want {
			t.Errorf("Got %d: wanted %d", got, want)
		}
	})

	t.Run("User may set bytes above the minimum", func(t *testing.T) {
		tm := setTokenManager(72)
		decodedStr, err := decodeStr(tm)
		if err != nil {
			t.Errorf("error decoding str %s", err.Error())
		}

		want := tm.BytesPerToken
		got := len(decodedStr)

		if got != want {
			t.Errorf("Got %d: wanted %d", got, want)
		}
	})
}

func setTokenManager(bytes int) *tokenManager {
	tm := &tokenManager{
		BytesPerToken: bytes,
	}
	return tm
}

func decodeStr(tm *tokenManager) (string, error) {
	tokenStr, err := tm.New()
	if err != nil {
		return "", fmt.Errorf("new Error: %w", err)
	}
	decodedStr, err := base64.URLEncoding.DecodeString(tokenStr)
	if err != nil {
		return "", fmt.Errorf("decode tokenstring error : %w", err)
	}

	return string(decodedStr), nil
}
