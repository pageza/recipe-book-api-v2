package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/middleware"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthMiddleware(t *testing.T) {
	// Use a known secret.
	secret := "testsecret"
	// Generate a token for a test user.
	token, err := utils.GenerateJWT("test-user-id", secret)
	assert.NoError(t, err)

	// Create a test Gin context.
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c.Request = req

	// Create the middleware.
	middlewareFunc := middleware.JWTAuth(secret)
	middlewareFunc(c)

	// After middleware executes, the userID should be set.
	userID, exists := c.Get("userID")
	assert.True(t, exists, "userID should be set in context")
	assert.Equal(t, "test-user-id", userID)
}
