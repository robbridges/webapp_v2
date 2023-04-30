package models

import (
	"testing"
)

const (
	// size can be different because of the hashing that we do on the token
	lengthAsString = 44
)

func TestTokenManager_New(t *testing.T) {
	// Test case 1: BytesPerToken is greater than or equal to the minimum value
	t.Run("Proper amount of bytes", func(t *testing.T) {
		tm := &tokenManager{
			BytesPerToken: 32,
		}

		tokenStr, err := tm.New()

		if err != nil {
			t.Errorf("Fail: error creating Token: %s", err.Error())
		}

		want := lengthAsString

		got := len(tokenStr)

		if got != want {
			t.Errorf("Got %d: wanted %d", got, want)
		}
	})
	t.Run("Improper amount of bytes", func(t *testing.T) {
		tm := &tokenManager{
			BytesPerToken: 6,
		}

		tokenStr, err := tm.New()

		if err != nil {
			t.Errorf("Fail: error creating Token: %s", err.Error())
		}

		want := lengthAsString

		got := len(tokenStr)

		if got != want {
			t.Errorf("Got %d: wanted %d", got, want)
		}
	})
}
