package routes

import (
	"Server/controller"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {

	r.GET("/user/getuser/:id", controller.GetUserById)
}
