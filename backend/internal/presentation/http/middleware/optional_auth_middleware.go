package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/application/port"
)

// OptionalAuthMiddleware attaches user identity when a valid Bearer token is
// present, but never rejects anonymous requests.
func OptionalAuthMiddleware(jwtService port.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		userID, role, err := jwtService.ValidateToken(parts[1])
		if err != nil {
			c.Next()
			return
		}

		c.Set(UserIDKey, userID)
		c.Set(UserRoleKey, role)
		c.Next()
	}
}
