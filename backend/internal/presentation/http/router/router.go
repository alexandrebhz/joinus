package router

import (
	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/presentation/http/handler"
	"github.com/startup-job-board/backend/internal/presentation/http/middleware"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/internal/domain/repository"
)

func NewRouter(
	authHandler *handler.AuthHandler,
	startupHandler *handler.StartupHandler,
	jobHandler *handler.JobHandler,
	fileHandler *handler.FileHandler,
	jwtService port.JWTService,
	startupRepo repository.StartupRepository,
	allowedOrigins []string,
) *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORSMiddleware(allowedOrigins))
	r.Use(middleware.LoggerMiddleware())

	// Health check
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	public := r.Group("/api/v1")
	{
		// Auth
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/refresh", authHandler.RefreshToken)

		// Public startups
		public.GET("/startups", startupHandler.List)
		public.GET("/startups/slug/:slug", startupHandler.GetBySlug)

		// Public jobs
		public.GET("/jobs", jobHandler.List)
	}

	// Protected routes (JWT)
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(jwtService))
	{
		// User
		protected.GET("/me", authHandler.GetMe)

		// Startups
		protected.POST("/startups", startupHandler.Create)
		protected.PUT("/startups/:id", startupHandler.Update)
		protected.GET("/startups/:id", startupHandler.Get)

		// Jobs
		protected.POST("/jobs", jobHandler.Create)
		protected.PUT("/jobs/:id", jobHandler.Update)
		protected.DELETE("/jobs/:id", jobHandler.Delete)

		// Files
		protected.POST("/upload", fileHandler.Upload)
	}

	// API Token routes
	tokenRoutes := r.Group("/api/v1/token")
	tokenRoutes.Use(middleware.APITokenMiddleware(startupRepo))
	{
		// Token-based startup operations
		tokenRoutes.GET("/startup", func(c *gin.Context) {
			startupID := middleware.GetStartupID(c)
			c.JSON(200, gin.H{"startup_id": startupID})
		})

		// Token-based job operations
		tokenRoutes.POST("/jobs", jobHandler.Create)
		tokenRoutes.PUT("/jobs/:id", jobHandler.Update)
		tokenRoutes.DELETE("/jobs/:id", jobHandler.Delete)
		tokenRoutes.GET("/jobs", jobHandler.List)
	}

	return r
}



