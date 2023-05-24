package models

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
)

type PostgressConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMODE  string
}

// Open will open a sql connection with the provided Postgres. Callers will need to ensure it's closed
func Open(config PostgressConfig) (*sql.DB, error) {
	db, err := sql.Open(
		"pgx",
		config.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("error Opening DB: %w", err)
	}
	return db, nil
}

func DefaultPostgresConfig() PostgressConfig {
	return PostgressConfig{
		Host:     viper.GetString("DATABASE_HOST"),
		Port:     viper.GetString("DATABASE_PORT"),
		User:     viper.GetString("DATABASE_USER"),
		Password: viper.GetString("DATABASE_PASSWORD"),
		Database: viper.GetString("DATABASE"),
		SSLMODE:  "disable",
	}
}

func DefaultPostgesTestConfig() PostgressConfig {
	return PostgressConfig{
		Host:     viper.GetString("TEST_DATABASE_HOST"),
		Port:     viper.GetString("TEST_DATABASE_PORT"), // Update the key to "TEST_DATABASE_PORT"
		User:     viper.GetString("TEST_DATABASE_USER"),
		Password: viper.GetString("TEST_DATABASE_PASSWORD"),
		Database: viper.GetString("TEST_DATABASE"),
		SSLMODE:  "disable",
	}
}

func (cfg PostgressConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMODE)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("Migrate: %v", err)
	}
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("Migrate: %v", err)
	}
	return nil
}
