package router

import (
	"github.com/gin-gonic/gin"
)

func registerApiRoutes(r *Router) {
	api := r.Router.Group("/api")
	{
		registerAnnouncementsRoutes(r, api)
	}
}

func registerAnnouncementsRoutes(r *Router, group *gin.RouterGroup) {
	announcements := group.Group("/announcements")
	{
		announcements.GET("/", r.Handler.GetAnnouncementList)
		announcements.POST("/:postId", r.Handler.CreateAnnouncement)
		announcements.PUT("/:postId", r.Handler.UpdateAnnouncement)
		announcements.DELETE("/:postId", r.Handler.DeleteAnnouncement)
	}
}
