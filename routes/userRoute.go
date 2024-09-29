package routes

import (
    "github.com/gin-gonic/gin"
    "chat-backend/controllers"
    "chat-backend/middleware"
)

func UserRoutes(router *gin.Engine) {
    // Public routes
    router.POST("/login", controllers.LoginUser)
    router.GET("/search", controllers.SearchUsers)

    // User routes group with JWT Middleware
    userGroup := router.Group("/user")
    userGroup.Use(middleware.JWTMiddleware()) // Applying JWT middleware to all routes under /user
    {
        userGroup.POST("/block", controllers.BlockUser)
		userGroup.POST("/register", controllers.RegisterUser)
        userGroup.GET("/verifyEmail/:token", middleware.VerifyEmail)
    }
}
