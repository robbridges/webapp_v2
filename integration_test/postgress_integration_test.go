package integration_test

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestOpen(t *testing.T) {
	db, err := setup(t)

	if err != nil {
		t.Errorf("Error opening db: %v", err)
	}
	defer deferDBClose(db, &err)
	defer teardown(t)

	if db == nil {
		t.Error("Database object should not be nil")
	}
	defer teardown(t)
}

func TestMigrate(t *testing.T) {

	db, err := setup(t)
	defer deferDBClose(db, &err)
	defer teardown(t)

	if err != nil {
		t.Errorf("Failed to drop table: %v", err)
	}

	if err != nil {
		t.Errorf("Error creating table: %v", err)
	}

	tableExists := tableExists(db, "users")
	if !tableExists {
		t.Errorf("expected users table to exist: %v", err)
	}

}

func tableExists(db *sql.DB, tableName string) bool {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = '%s')", tableName)
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}
