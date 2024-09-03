package database

import (
	"database/sql"
	"log"

	"chat-backend/config"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func Connect() (*Database, error) {
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to database")
	return &Database{DB: db}, nil
}
