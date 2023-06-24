package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

var (
	ErrEmailTaken = errors.New("models: email address already in use")
	ErrNoData     = errors.New("create sql: no rows in result set")
)

type FileError struct {
	Issue string
}

func (fe FileError) Error() string {
	return fmt.Sprintf("invalid file: %v", fe.Issue)
}

func CheckContentType(r io.ReadSeeker, allowedTypes []string) error {
	testBytes := make([]byte, 512)
	_, err := r.Read(testBytes)
	if err != nil {
		return fmt.Errorf("checking content type: %v", err)
	}
	_, err = r.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("checking content type: %v", err)
	}
	contentTypes := http.DetectContentType(testBytes)
	for _, t := range allowedTypes {
		if contentTypes == t {
			return nil
		}
	}
	return FileError{
		Issue: fmt.Sprintf("invalid content type"),
	}
}

func checkExtension(filename string, allowedExtensions []string) error {
	if hasExtension(filename, allowedExtensions) {
		return nil
	}
	return FileError{
		Issue: fmt.Sprintf("invalid extension: %v", filepath.Ext(filename)),
	}
}
