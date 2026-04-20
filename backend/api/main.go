package main

import (
	"Server/database"
	middlewares "Server/middleware"
	pb "Server/protos"
	"Server/routes"
	"Server/servergrpc"
	_ "Server/utils"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	err := database.Connect()
	if err != nil {
		panic(err)
	}

	//注册grpc服务
	listen, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("fail to listen:%v", err)
	}

	gprcServer := grpc.NewServer()
	pb.RegisterRealtimeChatServiceServer(gprcServer, &servergrpc.Server{})
	reflection.Register(gprcServer)

	log.Println("gprc running on port 5001")
	go func() {
		if err := gprcServer.Serve(listen); err != nil {
			log.Fatalf("fail to server grpc")
		}
	}()

	//end

	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	routes.SetupRoutes(r)

	r.Run() // listen and serve on
}
