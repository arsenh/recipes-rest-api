package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	DatabaseName string
}

// getEnv returns env var or fallback (if not set or empty)
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// mustGetEnv returns env var or fatal error
func mustGetEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	log.Fatalf("missing required environment variable: %s", key)
	return ""
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment / defaults")
	}

	cfg := &Config{
		Port:         getEnv("SERVER_PORT", ":8080"),
		DatabaseURL:  mustGetEnv("DATABASE_URL"),
		DatabaseName: getEnv("DATABASE_NAME", "test"),
	}

	//additional check
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("database URL is empty")
	}

	return cfg, nil
}
