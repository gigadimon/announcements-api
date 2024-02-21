package router

import (
	"announce-api/handlers"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Handler *handlers.Handler
	Router  *gin.Engine
}

func newRouter(handler *handlers.Handler) *Router {
	return &Router{
		Router:  gin.Default(),
		Handler: handler,
	}
}

func Init(handler *handlers.Handler) *gin.Engine {
	r := newRouter(handler)
	r.Router.Use(handler.CORSMiddleware())

	registerAuthRoutes(r)

	registerApiRoutes(r)

	return r.Router
}
