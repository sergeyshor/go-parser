package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// PG is a struct for storing Postgres connection settings
type PG struct {
	PostgresUser     string
	PostgresPassword string
	PostgresName     string
	PostgresPort     string
}

// Config is a struct for storing all required configuration parameters
type Config struct {
	*PG
	URL string
}

// New returns application config
func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// Postgres
	user, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		return nil, errors.New("POSTGRES_USER is not set")
	}

	pgPassword, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		return nil, errors.New("POSTGRES_PASSWORD is not set")
	}

	name, ok := os.LookupEnv("POSTGRES_NAME")
	if !ok {
		return nil, errors.New("POSTGRES_NAME is not set")
	}

	pgPort, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		return nil, errors.New("POSTGRES_PORT is not set")
	}

	// URL to parse
	parseUrl, ok := os.LookupEnv("PARSE_URL")
	if !ok {
		return nil, errors.New("PARSE_URL is not set")
	}

	return &Config{
		PG: &PG{
			PostgresUser:     user,
			PostgresPassword: pgPassword,
			PostgresName:     name,
			PostgresPort:     pgPort,
		},
		URL: parseUrl,
	}, nil
}
