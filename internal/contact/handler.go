package contact

import (
	"chat-backend/internal/database"
	"database/sql"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	db *sql.DB
}

type Contact struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func GetContactByID(c *gin.Context) {
	id := c.Param("id")
	var contact Contact

	if err := database.DB.First(&contact, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve contact"})
		}
		return
	}

	c.JSON(http.StatusOK, contact)
}
