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
	announcements.Use(r.Handler.Authenticate)
	{
		announcements.GET("/", r.Handler.GetAnnouncementList)
		announcements.GET("/:postId", r.Handler.GetAnnouncementById)
		announcements.POST("/", r.Handler.CreateAnnouncement)
		announcements.PATCH("/:postId", r.Handler.UpdateAnnouncement)
		announcements.DELETE("/:postId", r.Handler.DeleteAnnouncement)
	}

}
