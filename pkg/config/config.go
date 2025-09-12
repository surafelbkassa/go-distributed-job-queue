package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds application configuration
type Config struct {
	// Redis configuration
	RedisURL string

	// Server configuration
	ServerHost string
	ServerPort int

	// Worker configuration
	WorkerConcurrency int

	// Job configuration
	DefaultMaxAttempts int
	DefaultTimeout     time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		RedisURL:           getEnv("REDIS_URL", "redis://localhost:6379"),
		ServerHost:         getEnv("SERVER_HOST", "localhost"),
		ServerPort:         getEnvInt("SERVER_PORT", 8080),
		WorkerConcurrency:  getEnvInt("WORKER_CONCURRENCY", 5),
		DefaultMaxAttempts: getEnvInt("DEFAULT_MAX_ATTEMPTS", 3),
		DefaultTimeout:     getEnvDuration("DEFAULT_TIMEOUT", "30s"),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an environment variable as an integer with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvDuration gets an environment variable as a duration with a default value
func getEnvDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return 30 * time.Second
}