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
	ownerUser := factories.UserFactory(func(user *model.User) {
		user.ID = "00000000-0000-0000-0000-000000000000"
		user.UserName = "owner"
		user.Email = "owner@example.com"
		user.PasswordHash = hashPassword
	})
	exampleUser := factories.UserFactory(func(u *model.User) {
		u.ID = "00000000-0000-0000-0000-000000000001"
		u.UserName = "example"
		u.Email = "example@example.com"
		u.PasswordHash = hashPassword
	})
	gormDB.Create(ownerUser)
	gormDB.Create(exampleUser)

	log.Println("Full seeding complete.")
}
