package main

import (
	"example.com/internal/db"
	"example.com/internal/routers"
	"log"
	"os"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	gormDB, err := db.Connect(dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	r := routers.SetupRouter(gormDB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
