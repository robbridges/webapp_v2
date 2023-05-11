package integration_tests

import (
	"github.com/robbridges/webapp_v2/models"
	"testing"
)

func TestOpen(t *testing.T) {
	cfg := models.DefaultPostgesTestConfig()
	db, err := models.Open(cfg)

	if db == nil {
		t.Error("Database object should not be nil")
	}

	if err != nil {
		t.Errorf("There should be no error opening the database: %v", err)
	}

	if db != nil {
		db.Close()
	}
}
