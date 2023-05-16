package integration_tests

import (
	"database/sql"
	"github.com/robbridges/webapp_v2/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggerMiddleware(t *testing.T) {
	// Set up test database

	db, pool, resource := setup()

	defer db.Close()

	// Initialize DBLogger with the test database connection
	logger := &models.DBLogger{
		DB: db,
	}

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	loggerMiddleware := models.LoggerMiddleware(logger)(handlerFunc)

	// Test normal request handling
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	loggerMiddleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Verify no error logs were stored
	errorLogsCount, err := getErrorLogsCount(db)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if errorLogsCount != 0 {
		t.Errorf("Expected error logs count: 0, but got %d", errorLogsCount)
	}

	// Test panic handling
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

	// Verify error log was stored
	errorLogsCount, err = getErrorLogsCount(db)
	if err != nil {
		t.Errorf("Unexpected error getting logs")
	}
	if errorLogsCount != 1 {
		t.Errorf("Expected error logs count: 1, but got %d", errorLogsCount)
	}
	teardown(pool, resource)
}

func getErrorLogsCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM logs").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
