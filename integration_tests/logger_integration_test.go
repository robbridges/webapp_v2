package integration_tests

import (
	"database/sql"
	"errors"
	"github.com/robbridges/webapp_v2/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoggerMiddleware(t *testing.T) {
	// Set up test database
	db, err := setup(t)

	defer deferDBClose(db, &err)
	defer teardown(t)

	logger := &models.DBLogger{
		DB: db,
	}

	err = waitForPing(db, 10*time.Second)
	if err != nil {
		t.Errorf("Database timeout: %v", err)
	}

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	loggerMiddleware := models.LoggerMiddleware(logger)(handlerFunc)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	loggerMiddleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	errorLogsCount, err := getErrorLogsCount(db)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if errorLogsCount != 0 {
		t.Errorf("Expected error logs count: 0, but got %d", errorLogsCount)
	}

	panickingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Something went wrong!")
	})
	loggerMiddleware = models.LoggerMiddleware(logger)(panickingHandler)

	req = httptest.NewRequest("GET", "/", nil)
	rr = httptest.NewRecorder()

	loggerMiddleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, but got %d", http.StatusInternalServerError, rr.Code)
	}

	errorLogsCount, err = getErrorLogsCount(db)
	if err != nil {
		t.Errorf("Unexpected error getting logs")
	}
	if errorLogsCount != 1 {
		t.Errorf("Expected error logs count: 1, but got %d", errorLogsCount)
	}

	// Test logger.Create()
	err = logger.Create(errors.New("test error"))
	if err != nil {
		t.Errorf("Failed to create log entry: %v", err)
	}

	// Verify error log was stored
	errorLogsCount, err = getErrorLogsCount(db)
	if err != nil {
		t.Errorf("Unexpected error getting logs")
	}
	if errorLogsCount != 2 {
		t.Errorf("Expected error logs count: 2, but got %d", errorLogsCount)
	}

}

func getErrorLogsCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM logs").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
