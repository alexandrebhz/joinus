package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/presentation/http/response"
	"github.com/startup-job-board/backend/pkg/errors"
)

// RateLimitSettings configures per-IP sliding-window limits.
type RateLimitSettings struct {
	Enabled      bool
	Window       time.Duration
	DefaultLimit int // authenticated / general
	ListLimit    int // public list scrape targets
	DetailLimit  int // public detail scrape targets
	TrustedLimit int // SSR internal key
}

type visitorBucket struct {
	mu    sync.Mutex
	times []time.Time
}

type rateLimiterStore struct {
	mu       sync.Mutex
	visitors map[string]*visitorBucket
	window   time.Duration
}

func newRateLimiterStore(window time.Duration) *rateLimiterStore {
	s := &rateLimiterStore{
		visitors: make(map[string]*visitorBucket),
		window:   window,
	}
	go s.cleanupLoop()
	return s
}

func (s *rateLimiterStore) cleanupLoop() {
	ticker := time.NewTicker(s.window)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		cutoff := time.Now().Add(-s.window * 2)
		for key, v := range s.visitors {
			v.mu.Lock()
			if len(v.times) == 0 || v.times[len(v.times)-1].Before(cutoff) {
				delete(s.visitors, key)
			}
			v.mu.Unlock()
		}
		s.mu.Unlock()
	}
}

func (s *rateLimiterStore) allow(key string, limit int) (bool, time.Duration) {
	now := time.Now()
	windowStart := now.Add(-s.window)

	s.mu.Lock()
	v, ok := s.visitors[key]
	if !ok {
		v = &visitorBucket{}
		s.visitors[key] = v
	}
	s.mu.Unlock()

	v.mu.Lock()
	defer v.mu.Unlock()

	// Drop timestamps outside the window.
	kept := v.times[:0]
	for _, t := range v.times {
		if t.After(windowStart) {
			kept = append(kept, t)
		}
	}
	v.times = kept

	if len(v.times) >= limit {
		retryAfter := s.window - now.Sub(v.times[0])
		if retryAfter < time.Second {
			retryAfter = time.Second
		}
		return false, retryAfter
	}

	v.times = append(v.times, now)
	return true, 0
}

// RateLimitMiddleware applies stricter limits to public jobs/startups scrape targets.
func RateLimitMiddleware(cfg RateLimitSettings) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) { c.Next() }
	}

	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}
	if cfg.DefaultLimit <= 0 {
		cfg.DefaultLimit = 100
	}
	if cfg.ListLimit <= 0 {
		cfg.ListLimit = 30
	}
	if cfg.DetailLimit <= 0 {
		cfg.DetailLimit = 60
	}
	if cfg.TrustedLimit <= 0 {
		cfg.TrustedLimit = 300
	}

	store := newRateLimiterStore(cfg.Window)

	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Exempt health and Stripe webhooks.
		if strings.HasSuffix(path, "/health") || strings.HasSuffix(path, "/billing/webhook") {
			c.Next()
			return
		}

		limit := cfg.DefaultLimit
		bucket := "default"

		if IsInternalTrusted(c) {
			limit = cfg.TrustedLimit
			bucket = "trusted"
		} else if GetUserID(c) != "" {
			limit = cfg.DefaultLimit
			bucket = "auth"
		} else if isPublicListPath(path) {
			limit = cfg.ListLimit
			bucket = "list"
		} else if isPublicDetailPath(path) {
			limit = cfg.DetailLimit
			bucket = "detail"
		}

		ip := clientIP(c)
		key := bucket + ":" + ip
		allowed, retryAfter := store.allow(key, limit)
		if !allowed {
			c.Header("Retry-After", strconv.Itoa(int(retryAfter.Seconds())+1))
			c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
			response.Error(c, http.StatusTooManyRequests, errors.ErrRateLimited)
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Next()
	}
}

func isPublicListPath(path string) bool {
	return path == "/api/v1/jobs" || path == "/api/v1/startups"
}

func isPublicDetailPath(path string) bool {
	return strings.HasPrefix(path, "/api/v1/jobs/") ||
		strings.HasPrefix(path, "/api/v1/startups/slug/")
}

func clientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}
	if xri := strings.TrimSpace(c.GetHeader("X-Real-IP")); xri != "" {
		return xri
	}
	return c.ClientIP()
}
