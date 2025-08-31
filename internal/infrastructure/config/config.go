package config

import (
	"os"
	"strconv"
)

type Config struct {
	Security SecurityConfig
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	URL      string
	Host     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Port     int
}

type SecurityConfig struct {
	CSRFSecret    string
	SessionSecret string
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
		Security: SecurityConfig{
			CSRFSecret:    getEnvOrDefault("CSRF_SECRET", "csrf-secret-key"),
			SessionSecret: getEnvOrDefault("SESSION_SECRET", "session-secret-key"),
		},
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
