package router

import (
	"announce-api/handlers"

	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/signin", handlers.SignIn)
		auth.POST("/signup", handlers.SignUp)
	}
}
