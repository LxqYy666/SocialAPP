package routes

import (
	"Server/controller"
	middlewares "Server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {

	r.GET("/user/getuser/:id", controller.GetUserById)

	r.PATCH("/user/update/:id", middlewares.AuthMiddleWare, controller.UpdateUser)

	r.PATCH("/user/:id/following", middlewares.AuthMiddleWare, controller.FollowUser)

	r.GET("/user/sug/user", middlewares.AuthMiddleWare, controller.SuggestedUsers)
}
