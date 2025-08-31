package main

import (
	"example.com/internal/db"
	"example.com/internal/web"
	"log"
	"os"
)

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	gormDB, err := db.Connect(dbUrl)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	server := web.NewServer(gormDB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on :%s", port)

	if err := server.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
