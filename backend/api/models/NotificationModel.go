package models

import "time"

type NotificationUser struct {
	Name    string `json:"name" bson:"name"`
	Avartar string `json:"avatar,omitempty" bson:"avatar,omitempty"`
}

type Notification struct {
	ID               string    `json:"id,omitempty" bson:"_id,omitempty"`
	Details          string    `json:"details" bson:"details"`
	MainUID          string    `json:"mainUserId" bson:"mainUserId"`
	TargetUID        string    `json:"targetUserId" bson:"targetUserId"`
	IsReaded         bool      `json:"isReaded" bson:"isReaded"`
	CreatedAt        time.Time `json:"createdAt" bson:"createdAt"`
	NotificationUser `json:"notificationUser" bson:"notificationUser"`
}
