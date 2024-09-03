package chat

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Chat struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SenderID   uint      `json:"sender_id"`
	ReceiverID uint      `json:"receiver_id"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}

type ChatHandler struct {
	db *database.database
}

func NewHandler(db *database.database) *ChatHandler {
	return &ChatHandler{db: db}
}

func (h *ChatHandler) CreateChat(c *gin.Context) {
	var chat Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, chat)
}

func (h *ChatHandler) GetChatByID(c *gin.Context) {
	id := c.Param("id")
	var chat Chat

	if err := h.db.First(&chat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	c.JSON(http.StatusOK, chat)
}

func (h *ChatHandler) UpdateChat(c *gin.Context) {
	id := c.Param("id")
	var chat Chat

	if err := h.db.First(&chat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.db.Save(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chat"})
		return
	}

	c.JSON(http.StatusOK, chat)
}

func (h *ChatHandler) DeleteChat(c *gin.Context) {
	id := c.Param("id")
	var chat Chat

	if err := h.db.Delete(&chat, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted"})
}

func (h *ChatHandler) EditMessage(c *gin.Context) {
	messageID := c.Param("id")
	var request struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID, _ := c.Get("userID")

	result, err := h.db.Exec("UPDATE messages SET content = $1, edited = true, updated_at = now() WHERE id = $2 AND sender_id = $3", request.Content, messageID, userID)
	if err != nil || result.RowsAffected() == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit message or message not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message updated successfully"})
}
