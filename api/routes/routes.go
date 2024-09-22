package apiRoutes

import (
	"exchange/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/hello", controllers.Hello)
	}
}
