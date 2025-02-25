package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/healthcheck", nil)
	w := httptest.NewRecorder()

	handlers.HealthHandler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK")
	assert.Equal(t, "OK", w.Body.String(), "Expected response body to be 'OK'")
}
