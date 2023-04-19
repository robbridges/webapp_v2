package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const SessionTokenBytes = 32

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)

	if err != nil {
		return nil, fmt.Errorf("bytes: %w", err)
	}
	// we actually want to error check first, if we get this far and there's not enough bytes, and no error
	// we have a problem.
	if nRead < n {
		return nil, fmt.Errorf("bytes: Did not read enough random bytes")
	}

	return b, nil
}

// String returns a string from a random byte slice that is created in
// Bytes N is the number of bytes being used to generate string
func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("string: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// SessionToken is a very simple wrapper that takes the hard coded const and returns our token string, it uses our helper
// method byte, to create the byte slice with crypto, it converts it a string, and this method uses a constant preset value

func SessionToken() (string, error) {
	return String(SessionTokenBytes)
}

// GenerateRandByteSlice is for our csrf byte slice
func GenerateRandByteSlice() []byte {
	byteSlice := make([]byte, 32)
	if _, err := rand.Read(byteSlice); err != nil {
		fmt.Errorf("error generating byteSlice: %v", err)
	}

	return byteSlice
}
