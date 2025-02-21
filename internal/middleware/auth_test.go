package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/middleware"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
	secret := "testsecret"
	token, err := utils.GenerateJWT("test-user-id", secret)
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c.Request = req

	// Invoke JWTAuth middleware.
	middleware.JWTAuth(secret)(c)

	// Check that userID was set in the context.
	userID, exists := c.Get("userID")
	assert.True(t, exists, "Expected userID to be set in context")
	assert.Equal(t, "test-user-id", userID)
}

func TestJWTAuthMiddleware_MissingHeader(t *testing.T) {
	secret := "testsecret"

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/profile", nil)
	// No Authorization header.
	c.Request = req

	middleware.JWTAuth(secret)(c)

	// Expect 401 since the header is missing.
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status 401 for missing header")
	body := w.Body.String()
	assert.Contains(t, body, "missing or invalid token", "Response should indicate missing token")
}

func TestJWTAuthMiddleware_InvalidPrefix(t *testing.T) {
	secret := "testsecret"
	token, err := utils.GenerateJWT("test-user-id", secret)
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/profile", nil)
	// Set header with wrong prefix.
	req.Header.Set("Authorization", "Token "+token)
	c.Request = req

	middleware.JWTAuth(secret)(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status 401 for invalid header prefix")
	body := w.Body.String()
	assert.Contains(t, body, "missing or invalid token", "Response should indicate invalid token header")
}

func TestJWTAuthMiddleware_InvalidToken(t *testing.T) {
	secret := "testsecret"
	invalidToken := "this-is-not-a-valid-token"

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+invalidToken)
	c.Request = req

	middleware.JWTAuth(secret)(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status 401 for invalid token")
	body := w.Body.String()
	assert.Contains(t, body, "invalid token", "Response should indicate token parsing failed")
}

func TestLoggerMiddleware_Output(t *testing.T) {
	// Capture the logger output by redirecting Gin's default writer.
	var buf bytes.Buffer
	gin.DefaultWriter = &buf
	// Reset default writer after test.
	defer func() {
		gin.DefaultWriter = nil
	}()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.Logger())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verify that the logger output contains the expected debug message.
	output := buf.String()
	assert.True(t, strings.Contains(output, "DEBUG: Logger - request method:"), "Expected logger output to contain debug message")
}
