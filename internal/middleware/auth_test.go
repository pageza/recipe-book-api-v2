package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/middleware"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// init sets the global logger once at package initialization.
func init() {
	middleware.Log = zap.NewNop()
}

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
	// Ensure the global logger is set before running the test.
	middleware.Log = zap.NewNop()

	secret := "testsecret"
	token, err := utils.GenerateJWT("test-user-id", secret)
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c.Request = req

	// Invoke the JWTAuth middleware.
	middleware.JWTAuth(secret)(c)

	// Check that userID was set in the context.
	userID, exists := c.Get("userID")
	assert.True(t, exists, "Expected userID to be set in context")
	assert.Equal(t, "test-user-id", userID)
}

func TestJWTAuthMiddleware_MissingHeader(t *testing.T) {
	middleware.Log = zap.NewNop()

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
	middleware.Log = zap.NewNop()

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
	middleware.Log = zap.NewNop()

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
	// Create a test observer capturing debug-level logs.
	core, observedLogs := observer.New(zap.DebugLevel)
	testLogger := zap.New(core)

	// Override the global logger for this test.
	middleware.Log = testLogger

	gin.SetMode(gin.TestMode)
	router := gin.New()
	// Attach the Logger middleware.
	router.Use(middleware.Logger())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Perform a test request.
	req := httptest.NewRequest("GET", "/test", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify that at least one log entry contains our expected message.
	found := false
	for _, entry := range observedLogs.All() {
		if strings.Contains(entry.Message, "HTTP Request") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected logger output to contain debug message, got logs:\n%v", observedLogs.All())
	}
}
