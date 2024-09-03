package main

import (
	"chat-backend/config"
	"chat-backend/internal/database"
	"chat-backend/internal/router"
	"log"
)

func main() {
	config.LoadConfig()

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Setup the router with the DB connection
	r := router.SetupRouter(db.DB, config.SecretKey)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server run failed:", err)
	}
}
