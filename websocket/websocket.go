package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"exchange/models"

	"github.com/gorilla/websocket"
)

var activeConnections = struct {
	sync.RWMutex
	connections []*websocket.Conn
}{}

func AddConnection(conn *websocket.Conn) {
	activeConnections.Lock()
	defer activeConnections.Unlock()
	activeConnections.connections = append(activeConnections.connections, conn)
}

func RemoveConnection(conn *websocket.Conn) {
	activeConnections.Lock()
	defer activeConnections.Unlock()

	for i, c := range activeConnections.connections {
		if c == conn {
			activeConnections.connections = append(activeConnections.connections[:i], activeConnections.connections[i+1:]...)
			break
		}
	}
}

func BroadcastMessage(messageType int, message []byte) {
	activeConnections.RLock()
	defer activeConnections.RUnlock()

	for _, conn := range activeConnections.connections {
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println("Error sending message:", err)
		}
	}
}

func UpdateOrderBook(newOrderBook models.OrderBook) {
	message, err := json.Marshal(newOrderBook)
	if err != nil {
		log.Println("Error marshaling order book:", err)
		return
	}

	BroadcastMessage(websocket.TextMessage, message)
}
