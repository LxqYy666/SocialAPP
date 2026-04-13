package controller

import (
	"Server/database"
	"Server/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {

	userSchema := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user models.CreateUser

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid request body",
			"detail": err.Error(),
		})
		return
	}

	checkUser := userSchema.FindOne(ctx, bson.D{{Key: "email", Value: user.Email}}).Decode(&user)
	if checkUser == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	newUser := models.User{
		Email:     user.Email,
		Password:  string(hashedPassword),
		Name:      user.FirstName + " " + user.LastName,
		Followers: make([]string, 0),
		Following: make([]string, 0),
	}

	result, err := userSchema.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": result.InsertedID,
	})

}
