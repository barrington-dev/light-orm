package config

import "os"

type PostgresConfig struct {
	Store
}

type PostgresTestingConfig struct {
	Store
}

func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Store{
			Database: os.Getenv("DB_DATABASE"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}
}

func NewPostgresTestingConfig() *PostgresTestingConfig {
	return &PostgresTestingConfig{
		Store{
			Database: os.Getenv("DB_TEST_DATABASE"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}
}
