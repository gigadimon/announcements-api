package router

func registerAuthRoutes(r *Router) {
	auth := r.Router.Group("/auth")
	{
		auth.POST("/signin", r.Handler.SignIn)
		auth.POST("/signup", r.Handler.SignUp)
	}
}
