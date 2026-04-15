package routes

import (
	controller "Server/controller"
	middlewares "Server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupNotificationRoutes(r *gin.Engine) {
	r.GET("/notification/get", middlewares.AuthMiddleWare, controller.GetNotificationByUserId)
	r.PATCH("/notification/markread", middlewares.AuthMiddleWare, controller.MarkNotificationRead)
}
