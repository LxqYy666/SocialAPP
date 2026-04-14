package controller

import (
	"Server/database"
	"Server/models"
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreatePost(c *gin.Context) {
	var post models.CreateOrUpdatePost
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var userSchema = database.DB.Collection("users")
	var postSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user models.User
	strID, exists := c.Get("userId")
	if !exists {
		c.JSON(400, gin.H{"error": "User ID not found"})
		return
	}
	objID, err := bson.ObjectIDFromHex(strID.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
	userResult := userSchema.FindOne(ctx, bson.M{"_id": objID})
	if userResult.Err() != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if err := userResult.Decode(&user); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode user data"})
		return
	}

	newPost := models.Post{
		Creator:      strID.(string),
		Title:        post.Title,
		Message:      post.Message,
		Name:         user.Name,
		SelectedFile: post.SelectedFile,
		Likes:        make([]string, 0),
		Comments:     make([]string, 0),
		CreatedAt:    time.Now(),
	}

	result, err := postSchema.InsertOne(ctx, newPost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create post"})
		return
	}
	c.JSON(201, gin.H{
		"message": "Post created successfully",
		"postId":  result.InsertedID,
	})

}

func GetPostById(c *gin.Context) {
	var postSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// fmt.Println("Post ID:", c.Param("id")) // Debugging line to check the post ID being received
	objId, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	postResult := postSchema.FindOne(ctx, bson.M{"_id": objId})
	if postResult.Err() != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	var post models.Post
	if err := postResult.Decode(&post); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode post data"})
		return
	}

	c.JSON(200, gin.H{
		"post": post,
	})
}

func UpdatePost(c *gin.Context) {
	var postSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	objId, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	var authPost models.Post
	err = postSchema.FindOne(ctx, bson.M{"_id": objId}).Decode(&authPost)
	if err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	} else {
		if authPost.Creator != c.GetString("userId") {
			c.JSON(403, gin.H{"error": "Forbidden: You are not the creator of this post"})
			return
		}
	}

	var post models.CreateOrUpdatePost
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err = postSchema.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": post})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Post updated successfully",
	})

}

func GetAllPosts(c *gin.Context) {
	var postSchema = database.DB.Collection("posts")
	var userSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userId := c.Query("id")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	var user models.User
	objId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
	err = userSchema.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	user.Following = append(user.Following, userId)

	filter := bson.M{"creator": bson.M{"$in": user.Following}}

	totalPosts, err := postSchema.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to count posts"})
		return
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})
	findOptions.SetSkip(int64((page - 1) * 2))
	findOptions.SetLimit(2)

	cursor, err := postSchema.Find(ctx, filter, findOptions)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch posts"})
		return
	}
	defer cursor.Close(ctx)

	var posts []models.Post
	if err := cursor.All(ctx, &posts); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode posts"})
		return
	}

	c.JSON(200, gin.H{
		"posts":         posts,
		"currentPage":   page,
		"numberOfPages": (totalPosts + 1) / 2,
	})

}

func GetPostsUsersBySearch(c *gin.Context) {
	var userSchema = database.DB.Collection("users")
	var postSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	searchQuery := c.Query("searchQuery")
	if searchQuery == "" {
		c.JSON(400, gin.H{"error": "Search query is required"})
		return
	}

	filterPosts := bson.M{"$or": []bson.M{
		{"title": bson.M{"$regex": searchQuery, "$options": "i"}},
		{"message": bson.M{"$regex": searchQuery, "$options": "i"}},
	}}

	filterUsers := bson.M{"$or": []bson.M{
		{"name": bson.M{"$regex": searchQuery, "$options": "i"}},
		{"email": bson.M{"$regex": searchQuery, "$options": "i"}},
	}}

	postCursor, err := postSchema.Find(ctx, filterPosts)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to search posts"})
		return
	}
	defer postCursor.Close(ctx)

	var posts []models.Post
	if err := postCursor.All(ctx, &posts); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode posts"})
		return
	}

	userCursor, err := userSchema.Find(ctx, filterUsers)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to search users"})
		return
	}
	defer userCursor.Close(ctx)

	var users []models.User
	if err := userCursor.All(ctx, &users); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode users"})
		return
	}

	c.JSON(200, gin.H{
		"posts": posts,
		"users": users,
	})

}
