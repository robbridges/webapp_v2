package integration_tests

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func listFiles(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Failed to read directory: %w", err)
	}

	fmt.Printf("Files in directory %s:\n", dir)
	for _, file := range files {
		fmt.Println(file.Name())
	}
	return nil
}

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
	// Call the setup function
	setup()

	// Run tests and get the exit code
	exitCode := m.Run()

	// Call the teardown function
	teardown()

	// Exit with the test result
	os.Exit(exitCode)
}
