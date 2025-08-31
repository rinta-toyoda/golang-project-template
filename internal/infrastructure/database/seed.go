package database

import (
	"context"
	"log"

	"github.com/google/uuid"

	"example.com/internal/domain/entity"
	"example.com/internal/domain/repository"
	"example.com/pkg/security"
)

func SeedDatabase(userRepo repository.UserRepository, hasher security.PasswordHasher) error {
	ctx := context.Background()

	// Check if users already exist
	if _, err := userRepo.FindByEmail(ctx, "admin@example.com"); err == nil {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	// Create admin user
	hashedPassword, err := hasher.Hash("admin123456")
	if err != nil {
		return err
	}

	adminUser := &entity.User{
		ID:           uuid.NewString(),
		UserName:     "admin",
		Email:        "admin@example.com",
		PasswordHash: hashedPassword,
	}

	if err := userRepo.Create(ctx, adminUser); err != nil {
		return err
	}

	// Create test user
	hashedPassword, err = hasher.Hash("testpass123")
	if err != nil {
		return err
	}

	testUser := &entity.User{
		ID:           uuid.NewString(),
		UserName:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	}

	if err := userRepo.Create(ctx, testUser); err != nil {
		return err
	}

	log.Println("Database seeded successfully")
	return nil
}
