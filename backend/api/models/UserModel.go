package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email" binding:"required,email"`
	Password  string             `json:"password" bson:"password" binding:"required,min=5"`
	ImageUrl  string             `json:"imageUrl" bson:"imageUrl"`
	Bio       string             `json:"bio" bson:"bio"`
	Followers []string           `json:"followers" bson:"followers"`
	Following []string           `json:"following" bson:"following"`
}

type CreateUser struct {
	Email     string `json:"email" bson:"email" binding:"required,email"`
	Password  string `json:"password" bson:"password" binding:"required,min=5"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
}
