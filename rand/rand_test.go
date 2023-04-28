package rand

import (
	"crypto/rand"
	"errors"
	"testing"
)

type errorReader struct{}

func (r *errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("mocked error")
}

func TestBytes_Error(t *testing.T) {
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
}
