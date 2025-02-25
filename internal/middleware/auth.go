/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"
)

// JWTAuth validates the JWT token and stores the extended claims in the context.
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println("DEBUG: Auth - Authorization header:", authHeader)
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			fmt.Println("DEBUG: Auth - missing or invalid token header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println("DEBUG: Auth - token string:", tokenStr)
		claims, err := utils.ParseJWT(tokenStr, secret)
		if err != nil {
			fmt.Println("DEBUG: Auth - token parsing failed with error:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		fmt.Println("DEBUG: Auth - token valid, claims:", claims)
		// Instead of setting just the userID, we now store the entire extended claims.
		c.Set("userClaims", claims)
		c.Next()
	}
}

// Logger is a simple logging middleware.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Write log output to gin.DefaultWriter so tests can capture it.
		fmt.Fprintln(gin.DefaultWriter, "DEBUG: Logger - request method:", c.Request.Method, "path:", c.Request.URL.Path)
		c.Next()
	}
}
