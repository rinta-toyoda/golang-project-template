package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	URL      string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	SecretKey string
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("PORT", "8080"),
			Env:  getEnvOrDefault("ENV", "development"),
		},
		Database: DatabaseConfig{
			URL:      os.Getenv("DATABASE_URL"),
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvIntOrDefault("DB_PORT", 5432),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", ""),
			DBName:   getEnvOrDefault("DB_NAME", "app_db"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey: os.Getenv("JWT_SECRET_KEY"),
		},
	}

	if cfg.JWT.SecretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY environment variable is required")
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
