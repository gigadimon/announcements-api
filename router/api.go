package router

import (
	"announce-api/handlers"

	"github.com/gin-gonic/gin"
)

func registerApiRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		registerAnnouncementsRoutes(api)
	}
}

func registerAnnouncementsRoutes(group *gin.RouterGroup) {
	announcements := group.Group("/announcements")
	{
		announcements.GET("/", handlers.GetAnnouncementList)
		announcements.POST("/:postId", handlers.CreateAnnouncement)
		announcements.PUT("/:postId", handlers.UpdateAnnouncement)
		announcements.DELETE("/:postId", handlers.DeleteAnnouncement)
	}
}
