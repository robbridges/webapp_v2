package models

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMockLogger_Create(t *testing.T) {
	ml := &MockLogger{
		ErrorLog: []error{errors.New("first error"), errors.New("second error")},
	}

	newError := errors.New("third error to add")
	if err := ml.Create(newError); err != nil {
		t.Errorf("Mock logger returned error")
	}

	errorTable := ml.ErrorLog
	got := len(errorTable)
	want := 3
	if got != want {
		t.Errorf("The error was not appended to the error log, error log got size %d, want error log size %d",
			got, want,
		)

	}
	if errorTable[2].Error() != newError.Error() {
		t.Errorf("Order is wrong on error log")
	}
}

func TestLoggerMiddleware(t *testing.T) {
	mockLogger := &MockLogger{
		ErrorLog: []error{},
	}
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	loggerMiddleware := LoggerMiddleware(mockLogger)(handlerFunc)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	loggerMiddleware.ServeHTTP(rr, req)

	if len(mockLogger.ErrorLog) != 0 {
		t.Errorf("Unexpected error log length: %d", len(mockLogger.ErrorLog))
	}

	// Test panic handling
	panickingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Something went wrong!")
	})
	loggerMiddleware = LoggerMiddleware(mockLogger)(panickingHandler)

	req = httptest.NewRequest("GET", "/", nil)
	rr = httptest.NewRecorder()

	loggerMiddleware.ServeHTTP(rr, req)

	if len(mockLogger.ErrorLog) != 1 {
		t.Errorf("Expected error log length: 1, but got %d", len(mockLogger.ErrorLog))
	}
}
