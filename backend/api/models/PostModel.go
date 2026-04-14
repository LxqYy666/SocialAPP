package models

import "time"

type Post struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	Creator      string    `json:"creator" bson:"creator"`
	Title        string    `json:"title" bson:"title" binding:"required"`
	Message      string    `json:"message" bson:"message" bingding:"required"`
	Name         string    `json:"name" bson:"name"`
	SelectedFile string    `json:"selectedFile" bson:"selectedFile"`
	Likes        []string  `json:"likes" bson:"likes"`
	Comments     []string  `json:"comments" bson:"comments"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
}
