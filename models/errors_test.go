package models

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"testing"
)

func TestFileError_Error(t *testing.T) {
	// Prepare test data
	issue := "test issue"
	expectedErrorMessage := fmt.Sprintf("invalid file: %v", issue)

	// Create a FileError instance
	fileErr := FileError{Issue: issue}

	// Call the Error method
	errMsg := fileErr.Error()

	// Assert that the returned error message matches the expected value
	if errMsg != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got: '%s'", expectedErrorMessage, errMsg)
	}
}

func TestCheckContentType(t *testing.T) {
	// Prepare test data
	allowedTypes := []string{"image/jpeg", "image/png"}
	validContent := []byte("valid test data")
	invalidContent := []byte("invalid test data")

	// Create a test reader for valid content
	validReader := bytes.NewReader(validContent)

	// Create a test reader for invalid content
	invalidReader := bytes.NewReader(invalidContent)

	// Call the function under test with valid content
	err := CheckContentType(validReader, allowedTypes)

	// Call the function under test with invalid content
	err = CheckContentType(invalidReader, allowedTypes)
	expectedErr := FileError{Issue: "invalid content type"}
	if !errors.As(err, &expectedErr) {
		t.Errorf("Expected FileError with issue 'invalid content type', got: %v", err)
	}
}

func TestCheckExtension(t *testing.T) {
	// Prepare test data
	allowedExtensions := []string{".jpg", ".png"}
	validFilename := "image.jpg"
	invalidFilename := "document.pdf"

	// Call the function under test with valid extension
	err := checkExtension(validFilename, allowedExtensions)
	if err != nil {
		t.Errorf("Expected no error for valid extension, got: %v", err)
	}

	// Call the function under test with invalid extension
	err = checkExtension(invalidFilename, allowedExtensions)
	expectedErr := FileError{Issue: fmt.Sprintf("invalid extension: %v", filepath.Ext(invalidFilename))}
	if !errors.As(err, &expectedErr) {
		t.Errorf("Expected FileError with issue 'invalid extension: .pdf', got: %v", err)
	}
}
