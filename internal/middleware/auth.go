/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"
	"go.uber.org/zap"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		zap.L().Debug("Auth - Authorization header", zap.String("authHeader", authHeader))
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			zap.L().Debug("Auth - missing or invalid token header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		zap.L().Debug("Auth - token string", zap.String("token", tokenStr))
		claims, err := utils.ParseJWT(tokenStr, secret)
		if err != nil {
			zap.L().Error("Auth - token parsing failed", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		zap.L().Debug("Auth - token valid", zap.Any("claims", claims))
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func AuthLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		zap.L().Debug("Auth Logger - request", zap.String("method", c.Request.Method), zap.String("path", c.Request.URL.Path))
		c.Next()
	}
}
