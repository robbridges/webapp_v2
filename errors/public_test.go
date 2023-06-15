package errors

import (
	"errors"
	"reflect"
	"testing"
)

var (
	err       = errors.New("Standard error")
	pubError  = Public(err, "public message")
	testError = pubError.(PublicError)
)

func TestPublic(t *testing.T) {
	if _, ok := pubError.(PublicError); !ok {
		t.Errorf("pubError is not of Public error Type")
	}
}

func TestPublicError_Error(t *testing.T) {
	if testError.Error() == "standard error" {
		t.Errorf("The error message should match")
	}
}

func TestPublicError_Public(t *testing.T) {
	if testError.Public() != "public message" {
		t.Errorf("The public message is wrong")
	}
}

func TestPublicError_Unwrap(t *testing.T) {
	if reflect.DeepEqual(testError.Unwrap(), err) {
		t.Errorf("The underlying error should match")
	}
}
