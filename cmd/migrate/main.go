package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	migrationsPath := "infra/migrations"
	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	cmd := exec.Command("migrate", "-path", migrationsPath, "-database", dbUrl, "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
