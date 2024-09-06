package router

import (
	"chat-backend/internal/chat"
	"chat-backend/internal/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB, secretKey string) *gin.Engine {
	r := gin.Default()

	authMiddleware := middleware.JWTAuthMiddleware(secretKey)

	chatHandler := chat.NewHandler(db) // Ensure the handler is properly initialized
	r.POST("/chat", authMiddleware, chatHandler.CreateChat)
	r.GET("/chat/:id", authMiddleware, chatHandler.GetChatByID)
	r.PUT("/chat/:id", authMiddleware, chatHandler.UpdateChat)
	r.DELETE("/chat/:id", authMiddleware, chatHandler.DeleteChat)

	return r
}
