package config

import (
	"os"

	"github.com/joho/godotenv"
)

const defaultServerPort = "8080"

type Config struct {
	ServerPort string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		ServerPort: envOrDefault("SERVER_PORT", defaultServerPort),
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
