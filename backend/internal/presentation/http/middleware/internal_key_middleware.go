package middleware

import (
	"crypto/subtle"

	"github.com/gin-gonic/gin"
)

const InternalTrustedKey = "internal_trusted"

// InternalKeyMiddleware marks requests that present a valid X-Internal-Key
// (used by Next.js SSR / sitemap). Empty configured key disables the check.
func InternalKeyMiddleware(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if expectedKey == "" {
			c.Next()
			return
		}
		provided := c.GetHeader("X-Internal-Key")
		if provided != "" &&
			subtle.ConstantTimeCompare([]byte(provided), []byte(expectedKey)) == 1 {
			c.Set(InternalTrustedKey, true)
		}
		c.Next()
	}
}

func IsInternalTrusted(c *gin.Context) bool {
	v, exists := c.Get(InternalTrustedKey)
	if !exists {
		return false
	}
	trusted, ok := v.(bool)
	return ok && trusted
}
