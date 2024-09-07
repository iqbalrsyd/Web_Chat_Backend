package internal

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *Database) *gin.Engine {
	router := gin.Default()

	// Public routes
	router.POST("/register", RegisterHandler(db))
	router.POST("/login", LoginHandler(db))

	// Protected routes
	protected := router.Group("/")
	protected.Use(JWTAuthMiddleware())
	{
		protected.POST("/chats", CreateChatHandler(db))
		protected.POST("/groups", CreateGroupHandler(db))
		protected.GET("/chats/:id", GetChatByIDHandler(db))
	}

	return router
}
