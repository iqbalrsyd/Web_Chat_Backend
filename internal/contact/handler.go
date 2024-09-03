package contact

import (
	"chat-backend/internal/database"
	"database/sql"
	"net/http"
)

type Contact struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func GetContactByID(c *gin.Context) {
	id := c.Param("id")
	var contact Contact

	// Menjalankan query SQL untuk mengambil satu kontak berdasarkan ID
	err := database.DB.QueryRow("SELECT id, name, phone, email FROM contacts WHERE id = $1", id).Scan(&contact.ID, &contact.Name, &contact.Phone, &contact.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve contact"})
		}
		return
	}

	// Mengembalikan hasil sebagai JSON
	c.JSON(http.StatusOK, contact)
}

func (h *ChatHandler) AddFriend(c *gin.Context) {
	var Request struct {
		Friend uint `json:"friend_id"`
	}

	if err := c.ShouldBindJSON(&Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_id, _ := c.Get("user_id")

	_, err := h.DB.Exec("INSERT INTO contacts (user_id, friend_id, status) VALUES ($1, $2, $3)", user_id, Request.Friend, "pending")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add friend"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request sent"})
}
