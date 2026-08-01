package router

import (
	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/infrastructure/config"
	"github.com/startup-job-board/backend/internal/presentation/http/handler"
	"github.com/startup-job-board/backend/internal/presentation/http/middleware"
)

type RouterDeps struct {
	AuthHandler    *handler.AuthHandler
	StartupHandler *handler.StartupHandler
	JobHandler     *handler.JobHandler
	FileHandler    *handler.FileHandler
	ContactHandler *handler.ContactHandler
	BillingHandler *handler.BillingHandler
	TeamHandler    *handler.TeamHandler
	AdminHandler   *handler.AdminHandler
	JWTService     port.JWTService
	StartupRepo    repository.StartupRepository
	AllowedOrigins []string
	RateLimit      config.RateLimitConfig
	InternalKey    string
}

func NewRouter(deps RouterDeps) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware(deps.AllowedOrigins))
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.InternalKeyMiddleware(deps.InternalKey))
	r.Use(middleware.OptionalAuthMiddleware(deps.JWTService))
	r.Use(middleware.RateLimitMiddleware(middleware.RateLimitSettings{
		Enabled:      deps.RateLimit.Enabled,
		Window:       deps.RateLimit.Window,
		DefaultLimit: deps.RateLimit.Requests,
		ListLimit:    deps.RateLimit.ListLimit,
		DetailLimit:  deps.RateLimit.DetailLimit,
		TrustedLimit: deps.RateLimit.TrustedLimit,
	}))

	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	public := r.Group("/api/v1")
	public.Use(middleware.PublicCacheMiddleware())
	{
		public.POST("/auth/register", deps.AuthHandler.Register)
		public.POST("/auth/login", deps.AuthHandler.Login)
		public.POST("/auth/refresh", deps.AuthHandler.RefreshToken)
		public.GET("/auth/oauth/:provider", deps.AuthHandler.OAuthStart)
		public.GET("/auth/oauth/:provider/callback", deps.AuthHandler.OAuthCallback)
		public.POST("/auth/oauth/exchange", deps.AuthHandler.ExchangeOAuthCode)

		public.GET("/startups", deps.StartupHandler.List)
		public.GET("/startups/slug/:slug", deps.StartupHandler.GetBySlug)
		public.GET("/jobs", deps.JobHandler.List)
		public.GET("/jobs/:id", deps.JobHandler.Get)
		public.POST("/contact", deps.ContactHandler.Create)
		public.POST("/billing/webhook", deps.BillingHandler.Webhook)
	}

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(deps.JWTService))
	{
		protected.GET("/me", deps.AuthHandler.GetMe)

		protected.POST("/startups", deps.StartupHandler.Create)
		protected.PUT("/startups/:id", deps.StartupHandler.Update)
		protected.GET("/startups/:id", deps.StartupHandler.Get)

		protected.POST("/jobs", deps.JobHandler.Create)
		protected.PUT("/jobs/:id", deps.JobHandler.Update)
		protected.DELETE("/jobs/:id", deps.JobHandler.Delete)

		protected.POST("/upload", deps.FileHandler.Upload)
		protected.POST("/billing/checkout", deps.BillingHandler.CreateCheckout)
		protected.GET("/billing/status", deps.BillingHandler.Status)

		// Teams
		protected.POST("/teams", deps.TeamHandler.Create)
		protected.GET("/teams", deps.TeamHandler.ListMine)
		protected.GET("/teams/:id", deps.TeamHandler.Get)
		protected.PATCH("/teams/:id", deps.TeamHandler.Update)
		protected.GET("/teams/:id/members", deps.TeamHandler.ListMembers)
		protected.POST("/teams/:id/invitations", deps.TeamHandler.Invite)
		protected.PATCH("/teams/:id/members/:userId", deps.TeamHandler.UpdateMember)
		protected.DELETE("/teams/:id/members/:userId", deps.TeamHandler.RemoveMember)
		protected.GET("/teams/:id/roles", deps.TeamHandler.ListRoles)
		protected.GET("/teams/:id/startups", deps.TeamHandler.ListStartups)
		protected.POST("/teams/:id/startups/:startupId", deps.TeamHandler.LinkStartup)
		protected.DELETE("/teams/:id/startups/:startupId", deps.TeamHandler.UnlinkStartup)
		protected.POST("/invitations/accept", deps.TeamHandler.AcceptInvitation)

		// Platform admin
		admin := protected.Group("/admin")
		admin.Use(middleware.RequirePlatformAdmin())
		{
			admin.GET("/users", deps.AdminHandler.ListUsers)
			admin.PATCH("/users/:id", deps.AdminHandler.UpdateUser)
			admin.GET("/teams", deps.AdminHandler.ListTeams)
			admin.POST("/startups", deps.AdminHandler.CreateStartup)
			admin.PUT("/startups/:id/team", deps.AdminHandler.LinkStartupTeam)
		}
	}

	tokenRoutes := r.Group("/api/v1/token")
	tokenRoutes.Use(middleware.APITokenMiddleware(deps.StartupRepo))
	{
		tokenRoutes.GET("/startup", func(c *gin.Context) {
			startupID := middleware.GetStartupID(c)
			c.JSON(200, gin.H{"startup_id": startupID})
		})
		tokenRoutes.POST("/jobs", deps.JobHandler.Create)
		tokenRoutes.PUT("/jobs/:id", deps.JobHandler.Update)
		tokenRoutes.DELETE("/jobs/:id", deps.JobHandler.Delete)
		tokenRoutes.GET("/jobs", deps.JobHandler.List)
	}

	return r
}
