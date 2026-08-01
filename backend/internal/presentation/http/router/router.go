package router

import (
	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/presentation/http/handler"
	"github.com/startup-job-board/backend/internal/presentation/http/middleware"
)

func NewRouter(
	authHandler *handler.AuthHandler,
	startupHandler *handler.StartupHandler,
	jobHandler *handler.JobHandler,
	fileHandler *handler.FileHandler,
	contactHandler *handler.ContactHandler,
	billingHandler *handler.BillingHandler,
	teamHandler *handler.TeamHandler,
	adminHandler *handler.AdminHandler,
	jwtService port.JWTService,
	startupRepo repository.StartupRepository,
	allowedOrigins []string,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware(allowedOrigins))
	r.Use(middleware.LoggerMiddleware())

	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	public := r.Group("/api/v1")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/refresh", authHandler.RefreshToken)
		public.GET("/auth/oauth/:provider", authHandler.OAuthStart)
		public.GET("/auth/oauth/:provider/callback", authHandler.OAuthCallback)
		public.POST("/auth/oauth/exchange", authHandler.ExchangeOAuthCode)

		public.GET("/startups", startupHandler.List)
		public.GET("/startups/slug/:slug", startupHandler.GetBySlug)
		public.GET("/jobs", jobHandler.List)
		public.GET("/jobs/:id", jobHandler.Get)
		public.POST("/contact", contactHandler.Create)
		public.POST("/billing/webhook", billingHandler.Webhook)
	}

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(jwtService))
	{
		protected.GET("/me", authHandler.GetMe)

		protected.POST("/startups", startupHandler.Create)
		protected.PUT("/startups/:id", startupHandler.Update)
		protected.GET("/startups/:id", startupHandler.Get)

		protected.POST("/jobs", jobHandler.Create)
		protected.PUT("/jobs/:id", jobHandler.Update)
		protected.DELETE("/jobs/:id", jobHandler.Delete)

		protected.POST("/upload", fileHandler.Upload)
		protected.POST("/billing/checkout", billingHandler.CreateCheckout)
		protected.GET("/billing/status", billingHandler.Status)

		// Teams
		protected.POST("/teams", teamHandler.Create)
		protected.GET("/teams", teamHandler.ListMine)
		protected.GET("/teams/:id", teamHandler.Get)
		protected.PATCH("/teams/:id", teamHandler.Update)
		protected.GET("/teams/:id/members", teamHandler.ListMembers)
		protected.POST("/teams/:id/invitations", teamHandler.Invite)
		protected.PATCH("/teams/:id/members/:userId", teamHandler.UpdateMember)
		protected.DELETE("/teams/:id/members/:userId", teamHandler.RemoveMember)
		protected.GET("/teams/:id/roles", teamHandler.ListRoles)
		protected.GET("/teams/:id/startups", teamHandler.ListStartups)
		protected.POST("/teams/:id/startups/:startupId", teamHandler.LinkStartup)
		protected.DELETE("/teams/:id/startups/:startupId", teamHandler.UnlinkStartup)
		protected.POST("/invitations/accept", teamHandler.AcceptInvitation)

		// Platform admin
		admin := protected.Group("/admin")
		admin.Use(middleware.RequirePlatformAdmin())
		{
			admin.GET("/users", adminHandler.ListUsers)
			admin.PATCH("/users/:id", adminHandler.UpdateUser)
			admin.GET("/teams", adminHandler.ListTeams)
			admin.POST("/startups", adminHandler.CreateStartup)
			admin.PUT("/startups/:id/team", adminHandler.LinkStartupTeam)
		}
	}

	tokenRoutes := r.Group("/api/v1/token")
	tokenRoutes.Use(middleware.APITokenMiddleware(startupRepo))
	{
		tokenRoutes.GET("/startup", func(c *gin.Context) {
			startupID := middleware.GetStartupID(c)
			c.JSON(200, gin.H{"startup_id": startupID})
		})
		tokenRoutes.POST("/jobs", jobHandler.Create)
		tokenRoutes.PUT("/jobs/:id", jobHandler.Update)
		tokenRoutes.DELETE("/jobs/:id", jobHandler.Delete)
		tokenRoutes.GET("/jobs", jobHandler.List)
	}

	return r
}
