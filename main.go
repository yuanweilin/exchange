package main

import (
	"exchange/database"
	"exchange/middleware"
	"exchange/services"

	apiRoutes "exchange/api/routes"
	websocketRoutes "exchange/websocket/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	r := gin.Default()

	r.Use(middleware.NewCORSConfig())

	go services.GetLatestOrderPrices()
	apiRoutes.RegisterAPIRoutes(r)
	websocketRoutes.RegisterWebSocketRoutes(r)

	r.Run("0.0.0.0:8082")
}
