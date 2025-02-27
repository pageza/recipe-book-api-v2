package users_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/handlers/users"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAndLoginHandler(t *testing.T) {
	// Set Gin to test mode.
	gin.SetMode(gin.TestMode)

	// Use the dummy service.
	dummySvc := &dummyUserService{}
	handler := users.NewUserHandler(dummySvc, "testsecret")

	// Set up Gin router for HTTP registration and login.
	router := gin.Default()
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	// Test Registration via HTTP.
	registerPayload := users.RegisterInput{
		Username:    "testuser",
		Email:       "testuser@example.com",
		Password:    "testpassword",
		Preferences: map[string]interface{}{"diet": "vegan"},
	}
	payloadBytes, err := json.Marshal(registerPayload)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Test Login via HTTP.
	loginPayload := map[string]string{
		"email":    "testuser@example.com",
		"password": "testpassword",
	}
	loginBytes, err := json.Marshal(loginPayload)
	assert.NoError(t, err)

	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(loginBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	_, tokenExists := response["token"]
	assert.True(t, tokenExists, "expected a token in the login response")
}
