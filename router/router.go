package router

import "github.com/gin-gonic/gin"

func Init() *gin.Engine {
	router := gin.Default()

	registerAuthRoutes(router)

	registerApiRoutes(router)

	return router
}
