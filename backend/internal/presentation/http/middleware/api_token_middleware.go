package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/presentation/http/response"
)

const StartupIDKey = "startup_id"

func APITokenMiddleware(startupRepo repository.StartupRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// Extract token (could be "Bearer sb_..." or just "sb_...")
		token := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}

		startup, err := startupRepo.FindByAPIToken(c.Request.Context(), token)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set(StartupIDKey, startup.ID)
		c.Next()
	}
}

func GetStartupID(c *gin.Context) string {
	startupID, exists := c.Get(StartupIDKey)
	if !exists {
		return ""
	}
	return startupID.(string)
}



