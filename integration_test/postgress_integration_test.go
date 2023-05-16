package integration_test

import (
	"github.com/robbridges/webapp_v2/models"
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	setup()
	cfg := models.DefaultPostgesTestConfig()
	db, err := models.Open(cfg)
	if err != nil {
		t.Errorf("Error opening db: %v", err)
	}
	defer db.Close()
	defer teardown()

	if db == nil {
		t.Error("Database object should not be nil")
	}

	// Retry connection for up to 10 seconds with 1-second intervals
	timeout := 10 * time.Second
	waitForPing(db, timeout)
}
