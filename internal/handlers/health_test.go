package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	// Set Gin to test mode.
	gin.SetMode(gin.TestMode)

	// Create a response recorder and Gin context.
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Create a GET test request to pass into the Gin context.
	ctx.Request = httptest.NewRequest("GET", "/healthcheck", nil)

	// Call the HealthHandler using the Gin context.
	handlers.HealthHandler(ctx)

	// Verify the response.
	assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK")
	assert.Equal(t, "OK", w.Body.String(), "Expected response body to be 'OK'")
}
