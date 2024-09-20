package middleware

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func NewCORSConfig() gin.HandlerFunc {
	config := cors.DefaultConfig()

	origins := os.Getenv("CORS_ALLOW_ORIGINS")
	if origins != "" {
		config.AllowOrigins = strings.Split(origins, ",")
	} else {
		config.AllowOrigins = []string{"*"}
	}

	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	return cors.New(config)
}
