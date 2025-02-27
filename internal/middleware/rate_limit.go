package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitMiddleware returns a gin middleware function for rate limiting.
func RateLimitMiddleware(fillInterval time.Duration, capacity int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(fillInterval), capacity)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		c.Next()
	}
}
