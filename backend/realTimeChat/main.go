package main

import (
	"log"
	"net/http"
	"realTimeChat/realtime"

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
	manager := realtime.NewConnectionManager(realtime.GetFriends)
	if manager == nil {
		log.Fatal("failed to create connection manager")
	}

	r.GET("/ws/:id", func(ctx *gin.Context) {
		userID := ctx.Param("id")
		if userID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
			return
		}

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}
		defer conn.Close()
		defer manager.RemoveConnection(userID)

		manager.AddConnection(userID, conn)

		for {
			var msg realtime.Message
			if err := conn.ReadJSON(&msg); err != nil {
				log.Printf("WebSocket read error for %s: %v", userID, err)
				break
			}

			if msg.Receiver == "" || msg.Content == "" {
				continue
			}

			msg.Sender = userID
			manager.SendToReceiver(msg)
		}


	})

	r.Run(":8081")

}
