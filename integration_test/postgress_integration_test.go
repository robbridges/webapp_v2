package integration_test

import (
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	db, err := setup(t)

	if err != nil {
		t.Errorf("Error opening db: %v", err)
	}
	defer db.Close()
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
