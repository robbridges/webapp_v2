package integration_test

import (
	"database/sql"
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

	err = dropTableIfExists(db)
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}

	if err != nil {
		t.Fatalf("Error creating table: %v", err)
	}

}

func dropTableIfExists(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users CASCADE")
	return err
}
