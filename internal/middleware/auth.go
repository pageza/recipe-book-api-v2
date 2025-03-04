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

// JWTAuth validates a JWT token and authorizes the request.
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		Log.Debug("Auth - Authorization header", zap.String("header", authHeader))
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			Log.Debug("Auth - missing or invalid token header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		Log.Debug("Auth - token string", zap.String("token", tokenStr))
		claims, err := utils.ParseJWT(tokenStr, secret)
		if err != nil {
			Log.Debug("Auth - token parsing failed", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		Log.Debug("Auth - token valid", zap.Any("claims", claims))
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
