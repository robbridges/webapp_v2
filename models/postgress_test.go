package models

import (
	"github.com/spf13/viper"
	"testing"
)

func TestPostgressConfig_String(t *testing.T) {
	expected := "host=localhost port=5432 user=johndoe password=secret dbname=mydb sslmode=disable"
	cfg := PostgressConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "johndoe",
		Password: "secret",
		Database: "mydb",
		SSLMODE:  "disable",
	}

	result := cfg.String()

	if result != expected {
		t.Errorf("Unexpected result. Got %s, expected %s", result, expected)
	}
}

func TestDefaultPostgresConfig(t *testing.T) {

	viper.SetConfigFile("../local.env")
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	testCases := []struct {
		name     string
		configFn func() PostgressConfig
		expected PostgressConfig
	}{
		{
			name:     "Default config",
			configFn: DefaultPostgresConfig,
			expected: PostgressConfig{
				Host:     viper.GetString("DATABASE_HOST"),
				Port:     viper.GetString("DATABASE_PORT"),
				User:     viper.GetString("DATABASE_USER"),
				Password: viper.GetString("DATABASE_PASSWORD"),
				Database: viper.GetString("DATABASE"),
				SSLMODE:  "disable",
			},
		},
		{
			name:     "Default test config",
			configFn: DefaultPostgesTestConfig,
			expected: PostgressConfig{
				Host:     viper.GetString("TEST_DATABASE_HOST"),
				Port:     viper.GetString("TEST_DATABASE_PORT"),
				User:     viper.GetString("TEST_DATABASE_USER"),
				Password: viper.GetString("TEST_DATABASE_PASSWORD"),
				Database: viper.GetString("TEST_DATABASE"),
				SSLMODE:  "disable",
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := tc.configFn()

			if cfg.Host != tc.expected.Host {
				t.Errorf("unexpected host value: got %q, want %q", cfg.Host, tc.expected.Host)
			}
			if cfg.Port != tc.expected.Port {
				t.Errorf("unexpected port value: got %q, want %q", cfg.Port, tc.expected.Port)
			}
			if cfg.User != tc.expected.User {
				t.Errorf("unexpected user value: got %q, want %q", cfg.User, tc.expected.User)
			}
			if cfg.Password != tc.expected.Password {
				t.Errorf("unexpected password value: got %q, want %q", cfg.Password, tc.expected.Password)
			}
			if cfg.Database != tc.expected.Database {
				t.Errorf("unexpected database value: got %q, want %q", cfg.Database, tc.expected.Database)
			}
			if cfg.SSLMODE != tc.expected.SSLMODE {
				t.Errorf("unexpected SSLMODE value: got %q, want %q", cfg.SSLMODE, tc.expected.SSLMODE)
			}
		})
	}
}
