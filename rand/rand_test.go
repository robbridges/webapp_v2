package rand

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"testing"
)

type errorReader struct{}

// mock rand reader to error test
func (r *errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("mocked error")
}

func isString(v interface{}) bool {
	_, ok := v.(string)
	return ok
}

func TestBytes(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		byteNumber := 6
		b, err := Bytes(byteNumber)
		if err != nil {
			t.Errorf("Error creating byte slice")
		}
		if len(b) != byteNumber {
			t.Errorf("Not enough bytes were created")
		}
	})
	t.Run("Sad path, read error", func(t *testing.T) {
		randReader := rand.Reader

		defer func() {
			rand.Reader = randReader
		}()

		// replace rand.Reader with an errorReader
		rand.Reader = &errorReader{}

		b, err := Bytes(10)

		if b != nil {
			t.Errorf("Bytes should return nil slice on error")
		}

		if err == nil {
			t.Errorf("Bytes should return an error")
		}

		expectedError := "bytes: mocked error"
		if err.Error() != expectedError {
			t.Errorf("Bytes returned unexpected error message. expected=%q, actual=%q", expectedError, err.Error())
		}
	})
}

func TestString(t *testing.T) {
	bytes := 5

	string, err := String(bytes)
	if err != nil {
		t.Errorf("String did not return a string")
	}
	decodedStr, err := base64.URLEncoding.DecodeString(string)
	if len(decodedStr) != bytes {
		t.Errorf("String not created with enough bytes")
	}
	if !isString(string) {
		t.Errorf("Byte slice not correctly encoded into string")
	}

}
