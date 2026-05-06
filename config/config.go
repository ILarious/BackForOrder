package config

import (
	"os"
	"strconv"

	"github.com/ILarious/BackForOrder/pkg/postgres"
	"github.com/joho/godotenv"
)

const defaultServerPort = "8080"

type Config struct {
	ServerPort string
	Postgres   postgres.Config
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		ServerPort: envOrDefault("SERVER_PORT", defaultServerPort),
		Postgres: postgres.Config{
			Host:         envOrDefault("POSTGRES_HOST", "127.0.0.1"),
			Port:         envIntOrDefault("POSTGRES_PORT", 5432),
			User:         os.Getenv("POSTGRES_USER"),
			Password:     os.Getenv("POSTGRES_PASSWORD"),
			Database:     os.Getenv("POSTGRES_DB"),
			SSLMode:      envOrDefault("POSTGRES_SSLMODE", "disable"),
			MaxOpenConns: envIntOrDefault("POSTGRES_MAX_OPEN_CONNS", 20),
			MaxIdleConns: envIntOrDefault("POSTGRES_MAX_IDLE_CONNS", 20),
		},
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func envIntOrDefault(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
