package controller

import (
	"Server/database"
	"Server/models"
	"context"
	"net/http"
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

func UpdateUser(c *gin.Context) {
	var userSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tokenUserId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if tokenUserId.(string) != c.Param("id") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	objId, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var updateData models.UpdateUser
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	update := bson.M{}
	if updateData.Name != "" {
		update["name"] = updateData.Name
	}
	if updateData.Bio != "" {
		update["bio"] = updateData.Bio
	}
	if updateData.ImageUrl != "" {
		update["imageUrl"] = updateData.ImageUrl
	}

	result, err := userSchema.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	var updatedUser models.User
	err = userSchema.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedUser)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode updated user data"})
		return
	}

	c.JSON(200, gin.H{
		"data": updatedUser,
	})

}
