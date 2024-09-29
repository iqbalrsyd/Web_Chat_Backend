package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    MongoURI         string
    JWTSecret        string
    GoogleClientID   string
    GoogleSecret     string
    GoogleCallbackURL string
}

// LoadConfig loads environment variables from .env or system environment
func LoadConfig() (*Config, error) {
    // Load .env file if it exists
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using environment variables")
    }

    config := &Config{
        MongoURI:          getEnv("MONGO_URI", "mongodb://localhost:27017"),
        JWTSecret:         getEnv("JWT_SECRET", "your-jwt-secret"),
        GoogleClientID:    getEnv("GOOGLE_CLIENT_ID", ""),
        GoogleSecret:      getEnv("GOOGLE_CLIENT_SECRET", ""),
        GoogleCallbackURL: getEnv("GOOGLE_CALLBACK_URL", "/user/auth/google/callback"),
    }

    return config, nil
}

// Helper function to get the environment variable or return a default value
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
