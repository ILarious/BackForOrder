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
	Kafka      KafkaConfig
	WorkerPool WorkerPoolConfig
}

type KafkaConfig struct {
	Brokers              []string
	OrderRequestTopic    string
	OrderResponseTopic   string
	OrderResponseGroupID string
}

type WorkerPoolConfig struct {
	Size int
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
		Kafka: KafkaConfig{
			Brokers:              []string{envOrDefault("KAFKA_BROKER", "localhost:9092")},
			OrderRequestTopic:    envOrDefault("KAFKA_ORDER_REQUEST_TOPIC", "vk-blogger-orders"),
			OrderResponseTopic:   envOrDefault("KAFKA_ORDER_RESPONSE_TOPIC", "vk-blogger-results"),
			OrderResponseGroupID: envOrDefault("KAFKA_ORDER_RESPONSE_GROUP_ID", "back-for-order"),
		},
		WorkerPool: WorkerPoolConfig{
			Size: envIntOrDefault("WORKER_POOL_SIZE", 4),
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
