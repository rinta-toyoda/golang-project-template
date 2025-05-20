package seed

import (
	"example.com/internal/db"
	"example.com/internal/db/factories"
	"example.com/internal/model"
	"example.com/internal/utils"
	"log"
	"os"
)

func Run() {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	gormDB, err := db.Connect(dbUrl)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Users
	hashPassword, err := utils.HashPassword("password")
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}
	ownerUser := factories.UserFactory(func(u *model.User) {
		u.ID = "00000000-0000-0000-0000-000000000000"
		u.Email = "owner@example.com"
		u.Phone = "+14151234567"
		u.PasswordHash = hashPassword
	})
	exampleUser := factories.UserFactory(func(u *model.User) {
		u.ID = "00000000-0000-0000-0000-000000000001"
		u.Email = "example@example.com"
		u.Phone = "+442071838750"
		u.PasswordHash = hashPassword
	})
	gormDB.Create(ownerUser)
	gormDB.Create(exampleUser)

	log.Println("Full seeding complete.")
}
