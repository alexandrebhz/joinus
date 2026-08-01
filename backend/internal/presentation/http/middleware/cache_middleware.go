package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// PublicCacheMiddleware sets short Cache-Control headers on successful public
// GET responses for jobs/startups endpoints.
func PublicCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Request.Method != "GET" {
			return
		}
		if c.Writer.Status() < 200 || c.Writer.Status() >= 300 {
			return
		}

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		if !isCacheablePublicPath(path) {
			return
		}

		// Skip caching for authenticated dashboard callers.
		if GetUserID(c) != "" {
			c.Header("Cache-Control", "private, no-store")
			return
		}

		c.Header("Cache-Control", "public, max-age=30, s-maxage=60")
	}
}

func isCacheablePublicPath(path string) bool {
	if path == "/api/v1/jobs" || path == "/api/v1/startups" {
		return true
	}
	return strings.HasPrefix(path, "/api/v1/jobs/") ||
		strings.HasPrefix(path, "/api/v1/startups/slug/")
}
