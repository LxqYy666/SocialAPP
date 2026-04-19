package realtime

import (
	"log"
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
	onlineFriends  map[string][]string
	getUserFriends func(string) <-chan []string
	lock           sync.Mutex
}

func NewConnectionManager(getUserFriends func(string) <-chan []string) *ConnectionManager {
	if getUserFriends != nil {
		return &ConnectionManager{
			connections:    make(map[string]*websocket.Conn),
			onlineFriends:  make(map[string][]string),
			getUserFriends: getUserFriends,
		}
	}

	return nil
}

func (cm *ConnectionManager) AddConnection(userId string, conn *websocket.Conn) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	cm.connections[userId] = conn
	cm.onlineFriends[userId] = []string{}

	//notify online friends
	for friendId := range cm.onlineFriends {
		if friendId != userId && cm.IsFriend(userId, friendId) {
			cm.onlineFriends[friendId] = append(cm.onlineFriends[friendId], userId)
			err := cm.connections[friendId].WriteJSON(map[string]any{
				"onlineFriends": cm.onlineFriends[friendId],
			})
			if err != nil {
				log.Printf("Error notify %s about %s : %v", friendId, userId, err.Error())
				return
			}
		}
	}

	//update new user online friends

	go func() {
		for friends := range cm.getUserFriends(userId) {
			if friends == nil {
				continue
			}
			for _, friendId := range friends {
				if cm.connections[friendId] != nil {
					cm.onlineFriends[userId] = append(cm.onlineFriends[userId], friendId)
					err := cm.connections[userId].WriteJSON(map[string]any{
						"onlineFriends": cm.onlineFriends[friendId],
					})
					if err != nil {
						log.Printf("Error notify %s about %s : %v", userId, friendId, err.Error())
						return
					}
				}
			}
		}
	}()

}

func (cm *ConnectionManager) RemoveConnection(userId string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	delete(cm.connections, userId)
	delete(cm.onlineFriends, userId)

	for friendId := range cm.onlineFriends {
		for i, id := range cm.onlineFriends[friendId] {
			if id == userId {
				cm.onlineFriends[friendId] = append(cm.onlineFriends[friendId][:i], cm.onlineFriends[friendId][i+1:]...)
				err := cm.connections[friendId].WriteJSON(map[string]any{
					"onlineFriends": cm.onlineFriends[friendId],
				})
				if err != nil {
					log.Printf("Error notify %s about %s : %v", friendId, userId, err.Error())
					return
				}
				break
			}
		}
	}
}

func (cm *ConnectionManager) SendToReceiver(msg Message) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	if conn, ok := cm.connections[msg.Receiver]; ok {
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error Sending message to %s : %v", msg.Receiver, err.Error())
		}
		//TODO: SAVE TO DB BY GRPC
	} else {
		log.Printf("Receiver %s not found", msg.Receiver)
	}
}

func (cm *ConnectionManager) IsFriend(userId, friendId string) bool {
	friends := <-cm.getUserFriends(userId)
	if friends == nil {
		return false
	}
	for _, v := range friends {
		if v == friendId {
			return true
		}
	}
	return false
}
