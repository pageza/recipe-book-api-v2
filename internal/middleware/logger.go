package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Log is a global logger instance that can be used throughout the codebase.
var Log *zap.Logger

// Init initializes the zap logger. Use zap.NewProduction() in production mode.
// For now, we are using zap.NewDevelopment() for human-friendly logs.
func Init() error {
	var err error
	Log, err = zap.NewDevelopment()
	if err != nil {
		return err
	}
	return nil
}

// Sync flushes any buffered log entries.
func Sync() error {
	return Log.Sync()
}

// Logger returns a Gin middleware that logs HTTP requests using zap.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request.
		c.Next()

		duration := time.Since(startTime)
		zap.L().Info("Incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
