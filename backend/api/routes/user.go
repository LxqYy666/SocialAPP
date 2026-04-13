package routes

import (
	"Server/controller"
	middlewares "Server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {

	r.GET("/user/getuser/:id", controller.GetUserById)

	r.PATCH("/user/update/:id", middlewares.AuthMiddleWare, controller.UpdateUser)
}
