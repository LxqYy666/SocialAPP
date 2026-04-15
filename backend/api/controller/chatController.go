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
