package main

import (
	"chat-backend/config"
	"chat-backend/internal"
	"log"
)

func main() {
	// Load settings
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize MongoDB connection
	db, err := internal.Connect(config.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Setup and run router
	router := internal.SetupRouter(db)
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
