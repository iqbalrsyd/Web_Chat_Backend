package handlers

import (
	"chat-backend/internal"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateGroupHandler handles group creation
func CreateGroupHandler(db *internal.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var group internal.Group
		if err := c.BindJSON(&group); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		group.ID = primitive.NewObjectID()
		if err := internal.CreateGroup(db, &group); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create group"})
			return
		}
		c.JSON(http.StatusOK, group)
	}
}
