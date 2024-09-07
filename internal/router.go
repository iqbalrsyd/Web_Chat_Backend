package internal

import (
	"chat-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *Database) *gin.Engine {
	router := gin.Default()

	// Public routes
	router.POST("/register", handlers.RegisterUserHandler(db))
	router.POST("/login", handlers.LoginHandler(db))

	// Protected routes
	protected := router.Group("/")
	protected.Use(JWTAuthMiddleware())
	{
		protected.POST("/chats", handlers.CreateChatHandler(db))
		protected.GET("/chats/:id", handlers.GetChatByIDHandler(db))

		// Group routes
		protected.POST("/groups", handlers.CreateGroupHandler(db))

		// Routes for settings
		protected.GET("/settings", handlers.GetSettingsHandler(db))       // Mendapatkan pengaturan pengguna
		protected.POST("/settings", handlers.UpdateSettingsHandler(db))   // Memperbarui pengaturan pengguna
		protected.DELETE("/settings", handlers.DeleteSettingsHandler(db)) // Menghapus atau mereset pengaturan pengguna
	}

	return router
}
