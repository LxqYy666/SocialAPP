package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	r := gin.Default()

	r.GET("/ws/:id", func(ctx *gin.Context) {
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}
		defer conn.Close()

		for {
			id := ctx.Param("id")
			log.Printf("Received: %s", "hello "+id)

			if err := conn.WriteMessage(websocket.TextMessage, []byte("hello "+id)); err != nil {
				log.Printf("Write error: %v", err)
				break
			}
		}
	})

	r.Run(":8081")

}
