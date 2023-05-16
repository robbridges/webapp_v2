package integration_tests

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/viper"
	"os"
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

func setup() {
	dir := "../"

	cmd := exec.Command("make", "test-setup")
	cmd.Dir = dir

	err := cmd.Run()

	if err != nil {
		panic(err)
	}
}

func teardown() {
	dir := "../"

	cmd := exec.Command("make", "test-teardown")
	cmd.Dir = dir
	err := cmd.Run()

	if err != nil {
		panic(err)
	}

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

func TestMain(m *testing.M) {

	loadConfig()

	// Run tests and get the exit code
	exitCode := m.Run()

	// Exit with the test result
	os.Exit(exitCode)

}
