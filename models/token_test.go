package models

import (
	"bytes"
	"crypto/rand"
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

func TestNewTokenWithError(t *testing.T) {
	tm := &tokenManager{
		BytesPerToken: 16,
	}

	mockReader := &bytes.Reader{}
	randReader := rand.Reader
	rand.Reader = mockReader
	defer func() {
		rand.Reader = randReader
	}()

	_, err := tm.New()
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
	if err.Error() != "create session token: string: bytes: EOF" {
		t.Errorf("Expected error create session token: string: bytes: EOF, but got %v", err)
	}
}

func setTokenManager(bytes int) *tokenManager {
	tm := &tokenManager{
		BytesPerToken: bytes,
	}
	return tm
}

// DecodeStr decodes the str from base 64 encoding, since our rand package generates a random string and encodes it in
// base 64 we need to decode it to check the length of the str in our tests.
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
