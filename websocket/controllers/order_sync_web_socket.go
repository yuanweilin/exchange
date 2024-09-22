package controllers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var activeConnections = struct {
	sync.RWMutex
	connections []*websocket.Conn
}{}

func addConnection(conn *websocket.Conn) {
	activeConnections.Lock()
	defer activeConnections.Unlock()
	activeConnections.connections = append(activeConnections.connections, conn)
}

func removeConnection(conn *websocket.Conn) {
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

func OrderSyncWebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		return
	}
	defer conn.Close()

	addConnection(conn)
	defer removeConnection(conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
	}
}
