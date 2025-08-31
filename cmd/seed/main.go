package main

import (
	"log"

	"example.com/internal/app"
	"example.com/internal/domain/repository"
	"example.com/internal/infrastructure/database"
	"example.com/pkg/security"
)

func main() {
	container, err := app.BuildContainer()
	if err != nil {
		log.Fatalf("Failed to build container: %v", err)
	}

	err = container.Invoke(func(
		userRepo repository.UserRepository,
		hasher security.PasswordHasher,
	) error {
		return database.SeedDatabase(userRepo, hasher)
	})

	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Database seeding completed")
}
