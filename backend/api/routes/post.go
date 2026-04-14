package routes

import (
	"Server/controller"
	middlewares "Server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupPostRoutes(r *gin.Engine) {
	r.POST("/post/create", middlewares.AuthMiddleWare, controller.CreatePost)
	r.GET("/post/get/:id", controller.GetPostById)
}
