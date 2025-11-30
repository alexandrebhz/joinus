package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port         string
	DatabaseURL  string
	BackendURL   string
	BackendToken string
	Environment  string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	return &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/crawler?sslmode=disable"),
		BackendURL:   getEnv("BACKEND_URL", "http://localhost:8080"),
		BackendToken: getEnv("BACKEND_TOKEN", ""),
		Environment:  getEnv("ENVIRONMENT", "development"),
	}, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

