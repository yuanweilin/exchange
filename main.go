package main

import (
	"exchange/database"
	"exchange/middleware"
	"exchange/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	r := gin.Default()

	r.Use(middleware.NewCORSConfig())

	routes.RegisterRoutes(r)
	r.Run("0.0.0.0:8081")
}
