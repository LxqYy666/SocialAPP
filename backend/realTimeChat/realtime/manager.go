package realtime

import (
	"log"
	"realTimeChat/servergrpc"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
}

type ConnectionManager struct {
	connections    map[string]*websocket.Conn
	getUserFriends func(string) ([]string, error)
	lock           sync.RWMutex
}

func NewConnectionManager(getUserFriends func(string) ([]string, error)) *ConnectionManager {
	if getUserFriends != nil {
		return &ConnectionManager{
			connections:    make(map[string]*websocket.Conn),
			getUserFriends: getUserFriends,
		}
	}

	return nil
}

func (cm *ConnectionManager) AddConnection(userId string, conn *websocket.Conn) {
	cm.lock.Lock()

	cm.connections[userId] = conn
	cm.lock.Unlock()

	onlineFriends := cm.GetOnlineFriends(userId)
	err := conn.WriteJSON(map[string]any{
		"onlineFriends": onlineFriends,
	})
	if err != nil {
		log.Printf("Error send online friends to %s: %v", userId, err)
	}

}

func (cm *ConnectionManager) RemoveConnection(userId string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	delete(cm.connections, userId)
}

func (cm *ConnectionManager) SendToReceiver(msg Message) {
	err := servergrpc.SendMessageClient(msg.Sender, msg.Receiver, msg.Content)
	if err != nil {
		log.Printf("Error save message via grpc from %s to %s: %v", msg.Sender, msg.Receiver, err)
		return
	}

	cm.lock.RLock()
	defer cm.lock.RUnlock()
	if conn, ok := cm.connections[msg.Receiver]; ok {
		err = conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error Sending message to %s : %v", msg.Receiver, err.Error())
		}
	}
}

func (cm *ConnectionManager) GetOnlineFriends(userId string) []string {
	friends, err := cm.getUserFriends(userId)
	if err != nil {
		log.Printf("Error get friends for %s: %v", userId, err)
		return []string{}
	}

	cm.lock.RLock()
	defer cm.lock.RUnlock()

	onlineFriends := make([]string, 0, len(friends))
	for _, friendId := range friends {
		if cm.connections[friendId] != nil {
			onlineFriends = append(onlineFriends, friendId)
		}
	}

	return onlineFriends
}
