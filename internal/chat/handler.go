package chat

import (
	"database/sql"
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
	db *sql.DB
}

func NewHandler(db *sql.DB) *ChatHandler {
	return &ChatHandler{db: db}
}

func (h *ChatHandler) CreateChat(c *gin.Context) {
	var chat Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.Exec(
		"INSERT INTO chats (sender_id, receiver_id, message, created_at) VALUES ($1, $2, $3, $4)",
		chat.SenderID, chat.ReceiverID, chat.Message, chat.CreatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, chat)
}

func (h *ChatHandler) GetChatByID(c *gin.Context) {
	id := c.Param("id")
	var chat Chat

	row := h.db.QueryRow("SELECT id, sender_id, receiver_id, message, created_at FROM chats WHERE id = $1", id)
	if err := row.Scan(&chat.ID, &chat.SenderID, &chat.ReceiverID, &chat.Message, &chat.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat"})
		}
		return
	}

	c.JSON(http.StatusOK, chat)
}

func (h *ChatHandler) UpdateChat(c *gin.Context) {
	id := c.Param("id")
	var chat Chat

	row := h.db.QueryRow("SELECT id, sender_id, receiver_id, message, created_at FROM chats WHERE id = $1", id)
	if err := row.Scan(&chat.ID, &chat.SenderID, &chat.ReceiverID, &chat.Message, &chat.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat"})
		return
	}

	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := h.db.Exec(
		"UPDATE chats SET sender_id = $1, receiver_id = $2, message = $3 WHERE id = $4",
		chat.SenderID, chat.ReceiverID, chat.Message, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chat"})
		return
	}

	c.JSON(http.StatusOK, chat)
}

func (h *ChatHandler) DeleteChat(c *gin.Context) {
	id := c.Param("id")

	_, err := h.db.Exec("DELETE FROM chats WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted"})
}
