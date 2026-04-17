package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NotificationUser struct {
	Name    string `json:"name" bson:"name"`
	Avartar string `json:"avatar,omitempty" bson:"avatar,omitempty"`
}

type Notification struct {
	ID               bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Details          string        `json:"details" bson:"details"`
	MainUID          string        `json:"mainUserId" bson:"mainUserId"`
	TargetUID        string        `json:"targetUserId" bson:"targetUserId"`
	Type             string        `json:"type" bson:"type"`
	IsReaded         bool          `json:"isReaded" bson:"isReaded"`
	CreatedAt        time.Time     `json:"createdAt" bson:"createdAt"`
	NotificationUser `json:"notificationUser" bson:"notificationUser"`
}
