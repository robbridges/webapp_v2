package integration_tests

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func getProjectRoot() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("Failed to get the executable path: %w", err)
	}

	rootDir := filepath.Dir(filepath.Dir(exePath))
	return rootDir, nil
}

func loadConfig() {
	viper.SetConfigFile("../local.env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("init: %w", err))
	}
}

func setup() {
	rootDir, err := getProjectRoot()
	if err != nil {
		log.Fatalf("Failed to get the current working directory: %v", err)
	}

	cmd := exec.Command("make", "test-setup")
	cmd.Dir = rootDir

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	if err != nil {
		log.Printf("Command output:\n%s", out.String())
		log.Fatalf("Failed to set up Docker: %v", err)
	}
}

func teardown() {
	rootDir, err := getProjectRoot()
	if err != nil {
		log.Fatalf("Failed to get the current working directory: %v", err)
	}

	cmd := exec.Command("make", "test-teardown")
	cmd.Dir = rootDir
	err = cmd.Run()

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err != nil {
		log.Printf("Command output:\n%s", out.String())
		log.Fatalf("Failed to set up Docker: %v", err)
	}
}

//func TestMain(m *testing.M) {
//	loadConfig()
//	// Call the setup function
//	setup()
//
//	// Run tests and get the exit code
//	exitCode := m.Run()
//
//	// Call the teardown function
//	teardown()
//
//	// Exit with the test result
//	os.Exit(exitCode)
//}
