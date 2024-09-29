package controllers

import (
    "context"
    "net/http"
    "time"

    "chat-backend/config"
    "chat-backend/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "github.com/gin-gonic/gin"
)

// SendMessage sends a message to a chat (individual or group)
func SendMessage(c *gin.Context) {
    var messageData struct {
        ChatID   primitive.ObjectID `json:"chat_id"`
        SenderID primitive.ObjectID `json:"sender_id"`
        Content  string             `json:"content"`
    }

    if err := c.ShouldBindJSON(&messageData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    message := models.Message{
        ChatID:    messageData.ChatID,
        SenderID:  messageData.SenderID,
        Content:   messageData.Content,
        CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
    }

    collection := config.DB.Collection("messages")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := collection.InsertOne(ctx, message)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
        return
    }

    message.ID = result.InsertedID.(primitive.ObjectID)
    c.JSON(http.StatusOK, message)
}

// GetMessages retrieves all messages from a chat
func GetMessages(c *gin.Context) {
    chatID, _ := primitive.ObjectIDFromHex(c.Param("chat_id"))

    collection := config.DB.Collection("messages")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var messages []models.Message
    cursor, err := collection.Find(ctx, bson.M{"chat_id": chatID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var message models.Message
        if err := cursor.Decode(&message); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding message"})
            return
        }
        messages = append(messages, message)
    }

    if err := cursor.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
        return
    }

    c.JSON(http.StatusOK, messages)
}
