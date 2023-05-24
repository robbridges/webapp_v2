package integration_test

import (
	"database/sql"
	"fmt"
	"github.com/robbridges/webapp_v2/models"
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	db, err := setup(t)

	if err != nil {
		t.Errorf("Error opening db: %v", err)
	}
	defer deferDBClose(db, &err)
	defer teardown()

	if db == nil {
		t.Error("Database object should not be nil")
	}

	// Retry connection for up to 10 seconds with 1-second intervals

	err = waitForPing(db, 10*time.Second)
	if err != nil {
		t.Errorf("Database timeout: %v", err)
	}
}

func TestMigrate(t *testing.T) {

	db, err := setup(t)
	defer deferDBClose(db, &err)
	defer teardown()

	err = waitForPing(db, 10*time.Second)

	err = dropTableIfExists(db)
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}

	err = models.Migrate(db, "../migrations")
	if err == nil {
		t.Fatalf("Expected error due to duplicate table: %v", err)
	}

	tableExists := checkIfTableExists(db, "users")
	if !tableExists {
		t.Error("Expected 'users' table to exist after migration, but it doesn't.")
	}
}

func dropTableIfExists(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users CASCADE")
	return err
}

func checkIfTableExists(db *sql.DB, tableName string) bool {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = '%s')", tableName)
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}
