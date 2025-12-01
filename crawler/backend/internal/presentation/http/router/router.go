package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/crawler/internal/presentation/http/handler"
	"time"
)

// SetupRouter sets up the HTTP router
func SetupRouter(siteHandler *handler.SiteHandler, crawlHandler *handler.CrawlHandler, crawlLogHandler *handler.CrawlLogHandler) *gin.Engine {
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		sites := api.Group("/sites")
		{
			sites.POST("", siteHandler.Create)
			sites.GET("", siteHandler.List)
			sites.GET("/:id", siteHandler.Get)
			sites.PUT("/:id", siteHandler.Update)
			sites.DELETE("/:id", siteHandler.Delete)
			sites.POST("/:id/crawl", crawlHandler.Execute)
			sites.GET("/:id/logs", crawlLogHandler.GetLogs)
			sites.GET("/:id/logs/latest", crawlLogHandler.GetLatestLog)
		}
	}

	return r
}

