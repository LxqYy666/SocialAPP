package routes

import (
	"Server/controller"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine) {
	r.POST("/user/signup", controller.SignUp)
	r.POST("/user/signin", controller.SignIn)
}
