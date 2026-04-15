package controller

import (
	"Server/database"
	"Server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func SendMsg(c *gin.Context) {
	var messageSchema = database.DB.Collection("messages")
	var unReadedMsgSchema = database.DB.Collection("unReadedMsgs")

	var msg models.SendMsg
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	senderId := c.GetString("userId")
	var newMsg = models.Message{
		Sender:   senderId,
		Receiver: msg.Receiver,
		Content:  msg.Content,
	}

	_, err := messageSchema.InsertOne(c, newMsg)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to send message"})
		return
	}

	var unReadedMsg models.UnReadedMsg
	filter := bson.M{"mainUserId": msg.Receiver, "otherUserId": senderId}
	update := bson.M{"$inc": bson.M{"numOfUnReadedMsg": 1}, "$set": bson.M{"isReaded": false}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	err = unReadedMsgSchema.FindOneAndUpdate(c, filter, update, opts).Decode(&unReadedMsg)
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(500, gin.H{"error": "Failed to update unread message count"})
		return
	}

	c.JSON(200, gin.H{"message": "Message sent successfully"})
}
