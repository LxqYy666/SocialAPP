package main

import (
	"Server/database"
	"Server/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	err := database.Connect()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	routes.SetupRoutes(r)

	r.Run() // listen and serve on
}
