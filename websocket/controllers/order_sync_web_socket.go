package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"exchange/services"
	"exchange/websocket"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func OrderSyncWebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		return
	}
	defer conn.Close()

	websocket.AddConnection(conn)
	defer websocket.RemoveConnection(conn)

	orderBook := services.FetchOrderBookFromDatabase()
	message, err := json.Marshal(orderBook)
	if err != nil {
		log.Println("Error marshaling order book:", err)
		return
	}
	if err := conn.WriteMessage(ws.TextMessage, message); err != nil {
		log.Println("Error sending initial order book:", err)
		return
	}

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
	}
}
