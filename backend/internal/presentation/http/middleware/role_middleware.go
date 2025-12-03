package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/presentation/http/response"
)

// RequireRole creates a middleware that requires the user to have one of the specified roles
func RequireRole(allowedRoles ...entity.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleStr := GetUserRole(c)
		if userRoleStr == "" {
			response.Forbidden(c)
			c.Abort()
			return
		}

		userRole := entity.UserRole(userRoleStr)
		hasPermission := false
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.Forbidden(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireCompanyRole creates a middleware that allows only admin or startup_owner roles
// This is used for company/startup management operations
func RequireCompanyRole() gin.HandlerFunc {
	return RequireRole(entity.UserRoleAdmin, entity.UserRoleStartupOwner)
}

