package main

import (
	"log"
	pb "realTimeNotification/servegrpc"
	"sync"

	"github.com/gorilla/websocket"
)

func main() {
	wsmutex := &sync.Mutex{}
	wsConnections := make(map[string]*websocket.Conn)

	if err := pb.StartGrpcServer(wsmutex, wsConnections); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}

	select {}
}
