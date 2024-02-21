package router

func registerApiRoutes(r *Router) {
	api := r.Router.Group("/api")
	api.Use(r.Handler.Authenticate)
	{
		api.GET("/feed", r.Handler.GetGlobalFeed)
		api.GET("/my-posts", r.Handler.GetAuthorsAnnouncementList)
		api.GET("/:postId", r.Handler.GetAnnouncementById)
		api.POST("/", r.Handler.CreateAnnouncement)

		apiWithOnlyByAuthorAccess := api.Group("/")
		apiWithOnlyByAuthorAccess.Use(r.Handler.IsUserAnnounceAuthor)
		{
			apiWithOnlyByAuthorAccess.PUT("/:postId", r.Handler.UpdateAnnouncementById)
			apiWithOnlyByAuthorAccess.DELETE("/:postId", r.Handler.DeleteAnnouncementById)
			apiWithOnlyByAuthorAccess.GET("/:postId/switch-visibility", r.Handler.SwitchAnnounceVisibilityById)

			announcePhotos := apiWithOnlyByAuthorAccess.Group("/:postId/photos")
			{
				announcePhotos.PATCH("/", r.Handler.UploadNewAnnouncePhotosById)
				announcePhotos.DELETE("/:photoName", r.Handler.DeleteAnnouncePhotoById)
			}

		}

	}

}
