package integration_tests

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"testing"
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

func TestMain(m *testing.M) {

	loadConfig()

	// Run tests and get the exit code
	exitCode := m.Run()

	// Exit with the test result
	os.Exit(exitCode)

}
