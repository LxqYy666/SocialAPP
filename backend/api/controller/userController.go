package controller

import (
	"Server/database"
	"Server/models"
	"context"
	"net/http"
	"slices"
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

func FollowUser(c *gin.Context) {
	var userSchema = database.DB.Collection("users")
	var notificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	firstUserId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	secondUserId := c.Param("id")

	if firstUserId.(string) == secondUserId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}

	objFirstUserId, err := bson.ObjectIDFromHex(firstUserId.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
	objSecondUserId, err := bson.ObjectIDFromHex(secondUserId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var firstUser models.User
	err = userSchema.FindOne(ctx, bson.M{"_id": objFirstUserId}).Decode(&firstUser)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	var secondUser models.User
	err = userSchema.FindOne(ctx, bson.M{"_id": objSecondUserId}).Decode(&secondUser)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if slices.Contains(firstUser.Following, secondUser.ID.Hex()) {
		firstUser.Following = slices.Delete(firstUser.Following, slices.Index(firstUser.Following, secondUser.ID.Hex()), slices.Index(firstUser.Following, secondUser.ID.Hex())+1)
		secondUser.Followers = slices.Delete(secondUser.Followers, slices.Index(secondUser.Followers, firstUser.ID.Hex()), slices.Index(secondUser.Followers, firstUser.ID.Hex())+1)
	} else {
		firstUser.Following = append(firstUser.Following, secondUser.ID.Hex())
		secondUser.Followers = append(secondUser.Followers, firstUser.ID.Hex())

		// Create a new notification for the followed user
		newNotification := models.Notification{
			MainUID:   secondUserId,
			TargetUID: firstUserId.(string),
			Details:   firstUser.Name + " started following you.",
			NotificationUser: models.NotificationUser{
				Name:    firstUser.Name,
				Avartar: firstUser.ImageUrl,
			},
			CreatedAt: time.Now(),
		}
		_, err = notificationSchema.InsertOne(ctx, newNotification)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create notification"})
			return
		}
	}

	_, err = userSchema.UpdateOne(ctx, bson.M{"_id": objFirstUserId}, bson.M{"$set": bson.M{"following": firstUser.Following}})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update following list"})
		return
	}
	_, err = userSchema.UpdateOne(ctx, bson.M{"_id": objSecondUserId}, bson.M{"$set": bson.M{"followers": secondUser.Followers}})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update followers list"})
		return
	}

	c.JSON(200, gin.H{
		"secondUser": secondUser,
		"firstUser":  firstUser,
	})
}

func SuggestedUsers(c *gin.Context) {
	var userSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	objUserId, err := bson.ObjectIDFromHex(userId.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var currentUser models.User
	err = userSchema.FindOne(ctx, bson.M{"_id": objUserId}).Decode(&currentUser)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	//获取关注列表中用户关注的人的ID
	followingIds := make([]string, 0)
	for _, followingId := range currentUser.Following {
		var followingUser models.User
		objFollowingId, err := bson.ObjectIDFromHex(followingId)
		if err != nil {
			continue
		}
		err = userSchema.FindOne(ctx, bson.M{"_id": objFollowingId}).Decode(&followingUser)
		if err == nil {
			for _, id := range followingUser.Following {
				if id != currentUser.ID.Hex() && !slices.Contains(currentUser.Following, id) {
					followingIds = append(followingIds, id)
				}
			}
		}
	}

	//从数据库中获取这些用户的信息
	suggestedUsers := make([]models.User, 0)
	for _, id := range followingIds {
		var user models.User
		objId, err := bson.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		err = userSchema.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err == nil {
			suggestedUsers = append(suggestedUsers, user)
		}
	}

	c.JSON(200, gin.H{
		"suggestedUsers": suggestedUsers,
	})

}

func DeleteUser(c *gin.Context) {
	var userSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	paramUserId := c.Param("id")
	if userId.(string) != paramUserId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	objUserId, err := bson.ObjectIDFromHex(userId.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	result, err := userSchema.DeleteOne(ctx, bson.M{"_id": objUserId})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete user"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, gin.H{
		"message": "User deleted successfully",
	})
}
