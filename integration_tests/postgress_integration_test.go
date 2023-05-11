package integration_tests

import (
	"github.com/robbridges/webapp_v2/models"
	"testing"
)

func init() {
	cfg := models.DefaultPostgesTestConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func TestOpen(t *testing.T) {
	// Get the default test config
	config := models.DefaultPostgesTestConfig()

	// Call the Open function with the test config
	db, err := models.Open(config)
	if err != nil {
		t.Fatalf("Failed to open the database connection: %v", err)
	}

	// Close the database connection when the test is done
	defer db.Close()

	// Test a simple query to check if the connection works
	rows, err := db.Query("SELECT 1")
	if err != nil {
		t.Fatalf("Failed to execute test query: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Fatal("Expected one row, got none")
	}

	var result int
	if err := rows.Scan(&result); err != nil {
		t.Fatalf("Failed to scan test query result: %v", err)
	}

	if result != 1 {
		t.Fatalf("Expected result to be 1, got %d", result)
	}
}
