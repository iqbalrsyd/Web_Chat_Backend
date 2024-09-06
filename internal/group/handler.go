package group

import (
	"chat-backend/internal/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateGroup(c *gin.Context) {
	var request struct {
		Name    string `json:"name"`
		Members []uint `json:"members"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Create group
	var groupID uint
	err := DB.QueryRow("INSERT INTO groups (name, created_by) VALUES ($1, $2) RETURNING id", request.Name, userID).Scan(&groupID)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	// Add members
	for _, memberID := range request.Members {
		_, err := DB.Exec("INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, $3)", groupID, memberID, "member")
		if err != nil {
			log.Printf("Failed to add member to group: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
			return
		}
	}

	// Add creator as admin
	_, err = database.DB.Exec("INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, $3)", groupID, userID, "admin")
	if err != nil {
		log.Printf("Failed to add admin to group: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group created successfully"})
}
