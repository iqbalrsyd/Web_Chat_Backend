package group

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ChatHandler) CreateGroup(c *gin.Context) {
	var request struct {
		Name    string `json:"name"`
		Members []uint `json:"members"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	var GroupID uint
	err := h.db.QueryRow("INSERT INTO groups (name,created_by) VALUES ($1, $2) RETURNING id", request.Name, userID).Scan(&GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	for _, member := range request.Members {
		_, err := h.db.Exec("INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, $3)", GroupID, memberID, "member")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
			return
		}
	}

	_, err := h.db.Exec("INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, $3)", GroupID, userID, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group created successfully"})
}
