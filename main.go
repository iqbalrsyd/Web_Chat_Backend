package main

import (
    "chat-backend/config"
    "chat-backend/routes"
    "github.com/gin-gonic/gin"
)

func main() {
    config, err := config.LoadConfig()
    if err != nil {
        panic("Failed to load configuration")
    }

    err = config.ConnectDB()
    if err != nil {
        panic("Failed to connect to the database")
    }

    router := gin.Default()

	protected := router.Group("/api")
    protected.Use(func(c *gin.Context) {
        middleware.JWTMiddleware(cfg, c.Writer, c.Request)
    })

    routes.GroupRoutes(router)
	routes.MessageRoutes(router)
	routes.NotificationRoutes(router)
	routes.UserRoutes(router)
	routes.chatRoutes(router)

    router.Run(":8080")
}