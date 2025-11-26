package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/internal/presentation/http/response"
)

const UserIDKey = "user_id"
const UserRoleKey = "user_role"

func AuthMiddleware(jwtService port.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		token := parts[1]
		userID, role, err := jwtService.ValidateToken(token)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set(UserIDKey, userID)
		c.Set(UserRoleKey, role)
		c.Next()
	}
}

func GetUserID(c *gin.Context) string {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return ""
	}
	return userID.(string)
}



