package routes

import (
	"Server/controller"
	middlewares "Server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupChatRoutes(r *gin.Engine) {
	r.POST("/chat/sendMsg", middlewares.AuthMiddleWare, controller.SendMsg)
	r.GET("/chat/getMsgByNums", middlewares.AuthMiddleWare, controller.GetMsgByNums)
}
