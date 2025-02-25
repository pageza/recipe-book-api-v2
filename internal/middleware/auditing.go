package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		// Log relevant information. You could enhance this to log user IDs or other details.
		log.Printf("AUDIT: %s %s - %d in %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
	}
}
