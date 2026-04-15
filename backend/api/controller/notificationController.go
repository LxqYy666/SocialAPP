package controller

import (
	"Server/database"
	"Server/models"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetNotificationByUserId(c *gin.Context) {
	var notificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userId := c.GetString("userId")

	cursor, err := notificationSchema.Find(ctx, bson.M{"mainUID": userId})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch notifications"})
		return
	}
	defer cursor.Close(ctx)

	var notifications []models.Notification
	for cursor.Next(ctx) {
		var notification models.Notification
		if err := cursor.Decode(&notification); err != nil {
			c.JSON(500, gin.H{"error": "Failed to decode notification data"})
			return
		}
		notifications = append(notifications, notification)
	}

	if len(notifications) == 0 {
		c.JSON(200, gin.H{"notifications": []models.Notification{}})
		return
	}

	c.JSON(200, gin.H{"notifications": notifications})
}

func MarkNotificationRead(c *gin.Context) {
	var notificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userId := c.GetString("userId")

	filter := bson.M{"mainUID": userId, "isRead": false}
	update := bson.M{"$set": bson.M{"isRead": true}}

	result, err := notificationSchema.UpdateMany(ctx, filter, update)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to mark notifications as read"})
		return
	}
	c.JSON(200, gin.H{"message": "Notifications marked as read", "modifiedCount": result.ModifiedCount})
}
