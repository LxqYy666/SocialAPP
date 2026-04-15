package models

type UnReadedMsg struct {
	ID               string `json:"id" bson:"_id,omitempty"`
	MainUserID       string `json:"mainUserId" bson:"mainUserId"`
	OtherUserID      string `json:"otherUserId" bson:"otherUserId"`
	NumOfUnReadedMsg int    `json:"numOfUnReadedMsg" bson:"numOfUnReadedMsg"`
	IsReaded         bool   `json:"isReaded" bson:"isReaded"`
}
