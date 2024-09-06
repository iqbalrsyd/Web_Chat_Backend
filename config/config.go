package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	JWTSecret   string
	DatabaseURL string
)

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("Cannot load .env file, using environment variables")
	}

	JWTSecret = getEnv("JWT_SECRET", "defaultsecret")
	DatabaseURL = getEnv("DATABASE_URL", "postgres://user:password@localhost/dbname?sslmode=disable")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
