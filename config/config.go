package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	RedisURL     string
	RedisDB      int
	ServerPort   string
	WorkerCount  int
	MaxRetries   int
	QueueName    string
	StatusPrefix string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379"),
		RedisDB:      getEnvAsInt("REDIS_DB", 0),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		WorkerCount:  getEnvAsInt("WORKER_COUNT", 5),
		MaxRetries:   getEnvAsInt("MAX_RETRIES", 3),
		QueueName:    getEnv("QUEUE_NAME", "jobs"),
		StatusPrefix: getEnv("STATUS_PREFIX", "job:status:"),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}