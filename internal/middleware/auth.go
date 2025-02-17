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
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Example debug logging for every request
		fmt.Println("DEBUG: Logger - request method:", c.Request.Method, "path:", c.Request.URL.Path)
		c.Next()
	}
}
