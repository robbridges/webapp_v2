package library

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/robbridges/webapp_v2/models"
	"github.com/spf13/viper"
	"log"
	"os/exec"
	"testing"
	"time"
)

func loadConfig() {
	viper.SetConfigFile("../local.env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("init: %w", err))
	}
}

func setup(t *testing.T) (*sql.DB, error) {
	dir := "../"

	cmd := exec.Command("make", "test-setup")
	cmd.Dir = dir

	err := cmd.Run()

	if err != nil {
		panic(err)
	}

	cfg := models.DefaultPostgesTestConfig()
	db, err := models.Open(cfg)
	err = waitForPing(db, 10*time.Second)
	if err != nil {
		t.Errorf("Database never responded")
	}
	err = migrateUp(t)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	return db, err
}

func teardown(t *testing.T) {
	migrateDown(t)

	dir := "../"

	cmd := exec.Command("make", "test-teardown")
	cmd.Dir = dir
	err := cmd.Run()

	if err != nil {
		panic(err)
	}

}

func migrateUp(t *testing.T) error {
	cfg := models.DefaultPostgesTestConfig()
	db, err := models.Open(cfg)
	if err != nil {

		return err
	}
	err = models.Migrate(db, "../migrations")
	if err != nil {
		return err
	}
	return nil
}

func migrateDown(t *testing.T) error {
	cfg := models.DefaultPostgesTestConfig()
	db, err := models.Open(cfg)
	if err != nil {

		return err
	}
	err = goose.Down(db, "../migrations")
	if err != nil {
		return err
	}
	return nil
}

func waitForPing(db *sql.DB, timeout time.Duration) error {
	startTime := time.Now()
	deadline := startTime.Add(timeout)

	for time.Now().Before(deadline) {
		err := db.Ping()
		if err == nil {
			return nil // Condition met
		}

		// Sleep for a short interval before retrying
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("timeout waiting for DB.Ping")
}

func closeDB(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// deferDBClose now returns a function which can be deferred.
func deferDBClose(db *sql.DB, existingErr *error) func() {
	return func() {
		closeErr := closeDB(db)
		if closeErr != nil {
			if *existingErr == nil {
				*existingErr = closeErr
			} else {
				fmt.Printf("Warning: Got an error when closing DB: %v\n", closeErr)
			}
		}
	}
}

func TestMain(m *testing.M) {
	loadConfig()

	// Run tests and get the exit code
	exitCode := m.Run()

	// Check the exit code and log an error message if non-zero
	if exitCode != 0 {
		log.Fatal("Tests failed")
	}

	// No error occurred, so the tests passed
	log.Println("All tests passed")
}
