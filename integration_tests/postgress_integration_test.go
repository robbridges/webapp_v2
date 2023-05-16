package integration_tests

import (
	"testing"
)

func TestOpen(t *testing.T) {
	db, pool, resource := setup()
	defer db.Close()
	defer teardown(pool, resource)
	if db == nil {
		t.Error("Database object should not be nil")
	}

}
