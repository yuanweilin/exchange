package websocketRoutes

import (
	"exchange/websocket/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterWebSocketRoutes(r *gin.Engine) {
	ws := r.Group("/ws")
	{
		ws.GET("/ordersync", controllers.OrderSyncWebSocketHandler)
	}
}
