package config

import (
	"log"
	"os"
)

type Config struct {
	HTTPPort  string
	DB_DSN    string
	RabbitURL string
	LogLevel  string
}

func Load() *Config {
	return &Config{
		HTTPPort:  getEnv("HTTP_PORT", "8080"),
		DB_DSN:    getEnv("DB_DSN", "postgres://postgres:postgres@localhost:5432/taskforge?sslmode=disable"),
		RabbitURL: getEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672/"),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, default_value string) string {
	val := os.Getenv(key)
	if val == "" {
		if default_value == "" {
			log.Fatalf("missing required env variable: %s", key)
		}
		return default_value
	}
	return val
}
