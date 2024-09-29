package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"chat-backend/models"
)

// CreateChat creates a new chat (individual or group)
func CreateChat(c *gin.Context) {
	db := c.MustGet("db").(*mongo.Database)

	var request struct {
		Type    string   `json:"type" binding:"required"`    // Chat type: "individual" or "group"
		Members []string `json:"members" binding:"required"` // List of member IDs
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Convert member IDs from string to ObjectID
	var memberIDs []primitive.ObjectID
	for _, id := range request.Members {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID format"})
			return
		}
		memberIDs = append(memberIDs, objectID)
	}

	// Prepare chat document
	chat := models.Chat{
		Type:      request.Type,
		Members:   memberIDs,
		CreatedAt: time.Now(),
	}

	// Insert the chat into the database
	collection := db.Collection("chats")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, chat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	chat.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusOK, gin.H{"message": "Chat created successfully", "chat": chat})
}

// GetChatByID retrieves a chat by its ID
func GetChatByID(c *gin.Context) {
	db := c.MustGet("db").(*mongo.Database)
	chatIDParam := c.Param("chatID")

	// Convert the chat ID from string to ObjectID
	chatID, err := primitive.ObjectIDFromHex(chatIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	collection := db.Collection("chats")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var chat models.Chat
	err = collection.FindOne(ctx, bson.M{"_id": chatID}).Decode(&chat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat": chat})
}
