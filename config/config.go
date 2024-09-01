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
		log.Println("Tidak dapat memuat file .env, menggunakan environment variable")
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
