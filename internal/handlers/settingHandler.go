package handlers

import (
	"chat-backend/internal"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetSettingsHandler handles the request to retrieve user settings
func GetSettingsHandler(db *internal.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(primitive.ObjectID)
		settings, err := internal.GetSettings(db, userID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Unable to retrieve settings"})
			return
		}
		c.JSON(200, settings)
	}
}

// UpdateSettingsHandler handles the request to update user settings
func UpdateSettingsHandler(db *internal.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(primitive.ObjectID)
		var newSettings internal.Settings
		if err := c.BindJSON(&newSettings); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		if err := internal.UpdateSettings(db, userID, newSettings); err != nil {
			c.JSON(500, gin.H{"error": "Unable to update settings"})
			return
		}
		c.JSON(200, gin.H{"message": "Settings updated successfully"})
	}
}

// DeleteSettingsHandler handles the request to reset user settings to default
func DeleteSettingsHandler(db *internal.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(primitive.ObjectID)

		if err := internal.DeleteSetting(db, userID); err != nil {
			c.JSON(500, gin.H{"error": "Unable to reset settings"})
			return
		}
		c.JSON(200, gin.H{"message": "Settings reset successfully"})
	}
}
