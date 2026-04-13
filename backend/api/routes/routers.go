package routes

import (
	"Server/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.POST("/user/signup", controller.SignUp)
	r.POST("/user/signin", controller.SignIn)

}
