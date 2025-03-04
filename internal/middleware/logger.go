/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Log is the global Zap logger instance used throughout middleware.
var Log *zap.Logger

// InitLogger initializes the global Zap logger.
func InitLogger() error {
	var err error
	Log, err = zap.NewDevelopment() // Replace with zap.NewProduction() in production.
	if err != nil {
		return err
	}
	return nil
}

// SyncLogger flushes any buffered log entries.
func SyncLogger() error {
	return Log.Sync()
}

// Logger logs HTTP requests using the global Zap logger.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		Log.Debug("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path))
		c.Next()
	}
}
