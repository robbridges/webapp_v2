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

	expected := PostgressConfig{
		Host:     viper.GetString("DATABASE_HOST"),
		Port:     viper.GetString("DATABASE_PORT"),
		User:     viper.GetString("DATABASE_USER"),
		Password: viper.GetString("DATABASE_PASSWORD"),
		Database: viper.GetString("DATABASE"),
		SSLMODE:  "disable",
	}

	cfg := DefaultPostgresConfig()

	if cfg.Host != expected.Host {
		t.Errorf("unexpected host value: got %q, want %q", cfg.Host, expected.Host)
	}
	if cfg.Port != expected.Port {
		t.Errorf("unexpected port value: got %q, want %q", cfg.Port, expected.Port)
	}
	if cfg.User != expected.User {
		t.Errorf("unexpected user value: got %q, want %q", cfg.User, expected.User)
	}
	if cfg.Password != expected.Password {
		t.Errorf("unexpected password value: got %q, want %q", cfg.Password, expected.Password)
	}
	if cfg.Database != expected.Database {
		t.Errorf("unexpected database value: got %q, want %q", cfg.Database, expected.Database)
	}
	if cfg.SSLMODE != expected.SSLMODE {
		t.Errorf("unexpected SSLMODE value: got %q, want %q", cfg.SSLMODE, expected.SSLMODE)
	}
}
