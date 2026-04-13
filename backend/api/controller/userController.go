package controller

import (
	"Server/database"
	"Server/models"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetUserById(c *gin.Context) {
	var userSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	objId, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	userResult := userSchema.FindOne(ctx, bson.M{"_id": objId})
	if userResult.Err() != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	var user models.User
	if err := userResult.Decode(&user); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode user data"})
		return
	}

	c.JSON(200, gin.H{
		"user":  user,
		"posts": []string{}, // Placeholder for user posts, implement post retrieval logic as needed
	})

}
