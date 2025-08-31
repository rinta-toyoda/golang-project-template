package main

import (
	"errors"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func main() {
	migrationsPath := "deployments/migrations/migrations"
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Validate DATABASE_URL format to mitigate gosec G204
	if err := validateDatabaseURL(dbURL); err != nil {
		log.Fatalf("Invalid DATABASE_URL: %v", err)
	}

	// #nosec G204 - DATABASE_URL is validated above
	cmd := exec.Command("migrate", "-path", migrationsPath, "-database", dbURL, "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

func validateDatabaseURL(dbURL string) error {
	// Basic validation that this looks like a valid database URL
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		return err
	}

	// Ensure it's a postgresql URL
	if parsedURL.Scheme != "postgres" && parsedURL.Scheme != "postgresql" {
		return errors.New("only postgres/postgresql URLs are supported")
	}

	// Basic sanity checks
	if parsedURL.Host == "" {
		return errors.New("database host is required")
	}

	if !strings.Contains(parsedURL.Path, "/") || len(parsedURL.Path) <= 1 {
		return errors.New("database name is required")
	}

	return nil
}
