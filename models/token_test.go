package models

import (
	"encoding/base64"
	"testing"
)

func TestTokenManager_New(t *testing.T) {
	t.Run("Proper amount of bytes", func(t *testing.T) {
		tm := &tokenManager{
			BytesPerToken: MinBytesPerToken,
		}

		tokenStr, err := tm.New()
		decodedStr, err := base64.URLEncoding.DecodeString(tokenStr)

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
		tm := &tokenManager{
			BytesPerToken: 6,
		}

		tokenStr, err := tm.New()

		decodedStr, err := base64.URLEncoding.DecodeString(tokenStr)

		if err != nil {
			t.Errorf("Fail: error creating Token: %s", err.Error())
		}

		want := MinBytesPerToken

		got := len(decodedStr)

		if got != want {
			t.Errorf("Got %d: wanted %d", got, want)
		}
	})
}
