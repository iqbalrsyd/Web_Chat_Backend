package handlers

import (
	"chat-backend/internal"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetChatByIDHandler(db *internal.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		chatID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(chatID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
			return
		}

		chat, err := internal.GetChatByID(db, objectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve chat"})
			return
		}
		c.JSON(http.StatusOK, chat)
	}
}
