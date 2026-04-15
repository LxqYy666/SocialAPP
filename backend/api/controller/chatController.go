package controller

import (
	"Server/database"
	"Server/models"
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func SendMsg(c *gin.Context) {
	var messageSchema = database.DB.Collection("messages")
	var unReadedMsgSchema = database.DB.Collection("unReadedMsgs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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

	_, err := messageSchema.InsertOne(ctx, newMsg)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to send message"})
		return
	}

	var unReadedMsg models.UnReadedMsg
	filter := bson.M{"mainUserId": msg.Receiver, "otherUserId": senderId}
	update := bson.M{"$inc": bson.M{"numOfUnReadedMsg": 1}, "$set": bson.M{"isReaded": false}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	err = unReadedMsgSchema.FindOneAndUpdate(ctx, filter, update, opts).Decode(&unReadedMsg)
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(500, gin.H{"error": "Failed to update unread message count"})
		return
	}

	c.JSON(200, gin.H{"message": "Message sent successfully"})
}

func GetMsgByNums(c *gin.Context) {
	var messageSchema = database.DB.Collection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	from, err := strconv.Atoi(c.Query("from"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid 'from' parameter"})
		return
	}

	firstuid := c.Query("firstuid")
	seconduid := c.Query("seconduid")

	senderFilter := bson.M{"sender": firstuid, "receiver": seconduid}
	receiverFilter := bson.M{"sender": seconduid, "receiver": firstuid}
	filter := bson.M{"$or": []bson.M{senderFilter, receiverFilter}}

	options := options.Find().SetSort(bson.M{"_id": -1}).SetSkip(int64(from * 2)).SetLimit(2)

	cursor, err := messageSchema.Find(ctx, filter, options)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve messages"})
		return
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err = cursor.All(ctx, &messages); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode messages"})
		return
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	if len(messages) == 0 {
		c.JSON(200, gin.H{"messages": []models.Message{}})
		return
	}

	c.JSON(200, gin.H{"messages": messages})
}

func GetUserUnReadedMsg(c *gin.Context) {
	var unReadedMsgSchema = database.DB.Collection("unReadedMsgs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userId := c.GetString("userId")
	cursor, err := unReadedMsgSchema.Find(ctx, bson.M{"mainUserId": userId, "isReaded": false})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve unread messages"})
		return
	}
	defer cursor.Close(ctx)

	var unReadedMsgs []models.UnReadedMsg
	if err = cursor.All(ctx, &unReadedMsgs); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode unread messages"})
		return
	}

	var totalUnReadedMsg int = 0
	for _, unReadedMsg := range unReadedMsgs {
		totalUnReadedMsg += unReadedMsg.NumOfUnReadedMsg
	}

	c.JSON(200, gin.H{"totalUnReadedMsg": totalUnReadedMsg, "unReadedMsgs": unReadedMsgs})
}

func MaskUnReadedMsg(c *gin.Context) {
	var unReadedMsgSchema = database.DB.Collection("unReadedMsgs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userId := c.GetString("userId")
	otherUserId := c.Query("otherUserId")
	if userId == "" || otherUserId == "" {
		c.JSON(400, gin.H{"error": "User ID and Other User ID are required"})
		return
	}

	filter := bson.M{"mainUserId": userId, "otherUserId": otherUserId}
	update := bson.M{"$set": bson.M{"numOfUnReadedMsg": 0, "isReaded": true}}
	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var updateDoc bson.M
	err := unReadedMsgSchema.FindOneAndUpdate(ctx, filter, update, options).Decode(&updateDoc)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to mark messages as read"})
		return
	}

	isMarked := updateDoc != nil

	c.JSON(200, gin.H{
		"message":  "Messages marked as read successfully",
		"isMarked": isMarked,
	})

}
