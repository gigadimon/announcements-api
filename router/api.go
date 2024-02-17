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

		announcementsWithOnlyByAuthorAccess := announcements.Group("/")
		announcementsWithOnlyByAuthorAccess.Use(r.Handler.IsUserAnnounceAuthor)
		{
			announcementsWithOnlyByAuthorAccess.PUT("/:postId", r.Handler.UpdateAnnouncementById)
			announcementsWithOnlyByAuthorAccess.DELETE("/:postId", r.Handler.DeleteAnnouncementById)
			announcementsWithOnlyByAuthorAccess.GET("/:postId/switch-visibility", r.Handler.SwitchAnnounceVisibilityById)

			announcePhotos := announcementsWithOnlyByAuthorAccess.Group("/:postId/photos")
			{
				announcePhotos.PATCH("/", r.Handler.UploadNewAnnouncePhotosById)
				announcePhotos.DELETE("/:photoName", r.Handler.DeleteAnnouncePhotoById)
			}

		}

	}

}
