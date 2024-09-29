package routes

import (
    "chat-backend/controllers"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"
)

func ChatRoutes(router *gin.Engine, db *mongo.Database) {
    chatGroup := router.Group("/chats")
    {
        chatGroup.POST("/", func(c *gin.Context) { controllers.CreateChat(c, db) })
        chatGroup.GET("/:chatID", func(c *gin.Context) { controllers.GetChatByID(c, db) })
    }
}
