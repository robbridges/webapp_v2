package integration_tests

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/robbridges/webapp_v2/models"
	"github.com/spf13/viper"
	"log"
	"os"
	"testing"
)

func loadConfig() {
	viper.SetConfigFile("../local.env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("init: %w", err))
	}
}

//func setup() {
//	dir := "../"
//
//	cmd := exec.Command("make", "test-setup")
//	cmd.Dir = dir
//
//	err := cmd.Run()
//
//	if err != nil {
//		panic(err)
//	}
//}
//
//func teardown() {
//	dir := "../"
//
//	cmd := exec.Command("make", "test-teardown")
//	cmd.Dir = dir
//	err := cmd.Run()
//
//	if err != nil {
//		panic(err)
//	}
//
//}

func setup() (*sql.DB, *dockertest.Pool, *dockertest.Resource) {
	cfg := models.DefaultPostgesTestConfig()
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Failed to create Docker pool: %v", err)
	}

	resource, err := pool.Run("postgres", "13-alpine", []string{
		"POSTGRES_USER=dev",
		"POSTGRES_PASSWORD=test",
		"POSTGRES_DB=testdb",
		"POSTGRES_PORT=5433",
	})
	if err != nil {
		log.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	conn, err := models.Open(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	return conn, pool, resource
}

func teardown(pool *dockertest.Pool, resource *dockertest.Resource) {
	if resource != nil {
		err := pool.Purge(resource)
		if err != nil {
			log.Fatalf("Failed to purge Docker resource: %v", err)
		}
	}
}

func TestMain(m *testing.M) {

	loadConfig()

	// Run tests and get the exit code
	exitCode := m.Run()

	// Exit with the test result
	os.Exit(exitCode)

}
