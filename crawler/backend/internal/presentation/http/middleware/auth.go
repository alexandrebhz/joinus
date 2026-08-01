package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIAuth returns middleware that validates API access via Bearer token or X-API-Key.
// When apiToken is empty the middleware is a no-op (development without a configured token).
func APIAuth(apiToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if apiToken == "" {
			c.Next()
			return
		}

		token := extractToken(c)
		if token == "" || token != apiToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	if apiKey := c.GetHeader("X-API-Key"); apiKey != "" {
		return apiKey
	}

	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}

	return ""
}
