package database

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func Connect() error {
	uri := "mongodb://localhost:27017"
	dbName := "social"

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return err
	}
	Client = client
	DB = Client.Database(dbName)
	fmt.Println("Connected to MongoDB successfully")
	return nil
}
