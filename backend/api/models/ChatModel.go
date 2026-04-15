package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Message struct {
	ID       bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Sender   string        `json:"sender" bson:"sender"`
	Receiver string        `json:"receiver" bson:"receiver"`
	Content  string        `json:"content" bson:"content" binding:"required"`
}

type SendMsg struct {
	Receiver string `json:"receiver" bson:"receiver" binding:"required"`
	Content  string `json:"content" bson:"content" binding:"required"`
}
