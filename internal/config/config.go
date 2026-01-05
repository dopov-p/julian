package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Env      string
}

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	DSN     string
	SSLMode string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			DSN:     getEnv("POSTGRES_PG_DSN", "postgres://db-dopov-p.julian-local:db-dopov-p.julian-local@localhost:5432/julian-local"),
			SSLMode: getEnv("POSTGRES_SSL_MODE", "disable"),
		},
		Env: getEnv("ENVIRONMENT", "development"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
