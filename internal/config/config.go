package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	HTTPPort      string
	DB_DSN        string
	RabbitURL     string
	LogLevel      string
	JWTSecret     string
	JWTExpiration time.Duration
}

func Load() *Config {
	config := &Config{
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
		DB_DSN:        getEnv("DB_DSN", "postgres://postgres:postgres@localhost:5432/taskforge?sslmode=disable"),
		RabbitURL:     getEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672/"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTExpiration: 24 * time.Hour,
	}

	if config.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	if len(config.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters long for security")
	}

	return config
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
