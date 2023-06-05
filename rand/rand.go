package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

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
