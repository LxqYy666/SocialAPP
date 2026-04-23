package realtime

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Notification struct {
	ID        string    `json:"id"`
	Details   string    `json:"details"`
	MainUID   string    `json:"mainuid"`
	TargetID  string    `json:"targetid"`
	Type      string    `json:"type"`
	IsReaded  bool      `json:"isreaded"`
	CreatedAt time.Time `json:"createdAt"`
	User      User      `json:"user"`
}

type User struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartWebSocketServer(wsmutex *sync.Mutex, wsConnections map[string]*websocket.Conn) {
	r := gin.Default()

	r.GET("/ws/:id", func(ctx *gin.Context) {
		userID := ctx.Param("id")
		// WebSocket connection logic will be handled here
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			return
		}
		wsmutex.Lock()
		wsConnections[userID] = conn
		wsmutex.Unlock()

		defer func() {
			wsmutex.Lock()
			delete(wsConnections, userID)
			wsmutex.Unlock()
			conn.Close()
		}()

		for {
			var notification Notification
			if err := conn.ReadJSON(&notification); err != nil {
				break
			}

			err := conn.WriteJSON(notification)
			if err != nil {
				break
			}
		}
	})

	r.Run(":8082")

}
