package models

import "go.mongodb.org/mongo-driver/v2/bson"

type UnReadedMsg struct {
	ID               bson.ObjectID `json:"id" bson:"_id,omitempty"`
	MainUserID       string        `json:"mainUserId" bson:"mainUserId"`
	OtherUserID      string        `json:"otherUserId" bson:"otherUserId"`
	NumOfUnReadedMsg int           `json:"numOfUnReadedMsg" bson:"numOfUnReadedMsg"`
	IsReaded         bool          `json:"isReaded" bson:"isReaded"`
}
