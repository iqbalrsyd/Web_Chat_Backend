package routes

import (
    "chat-backend/controllers"
    "github.com/gin-gonic/gin"
)

func MessageRoutes(router *gin.Engine) {
    message := router.Group("/messages")
    {
        message.POST("/", controllers.SendMessage)
        message.GET("/:chat_id", controllers.GetMessages)
    }
}
