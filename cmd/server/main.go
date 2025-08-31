package main

import (
	"log"

	"example.com/internal/app"
	"example.com/internal/infrastructure/database"
	"gorm.io/gorm"
)

func main() {
	container, err := app.BuildContainer()
	if err != nil {
		log.Fatalf("Failed to build container: %v", err)
	}

	// Run database migrations
	if err := container.Invoke(func(db *gorm.DB) error {
		return database.Migrate(db)
	}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	server, err := app.NewServer(container)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
