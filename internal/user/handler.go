package user

import (
	"chat-backend/internal/auth"
	"chat-backend/internal/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Example to clear password after hashing for security
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	u.PasswordHash = string(passwordHash)
	u.Password = "" // Clear the plaintext password

	// Insert user into database
	err = database.DB.QueryRow(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
		u.Username, u.PasswordHash,
	).Scan(&u.ID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			http.Error(w, "Username already exists", http.StatusConflict)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(strconv.Itoa(u.ID))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
