package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	"github.com/pageza/recipe-book-api-v2/internal/middleware"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestRouter configures an in-memory test DB, sets up the repository, service, handlers,
// and configures a Gin router with the /register, /login, and /profile endpoints.
func setupTestRouter() *gin.Engine {
	// Open an in-memory SQLite database.
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to in-memory sqlite")
	}

	// AutoMigrate will create the users table.
	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic(err)
	}

	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	handler := handlers.NewUserHandler(svc, "testsecret")

	router := gin.Default()
	// Register endpoints.
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)
	// Protect /profile using the JWT middleware.
	router.GET("/profile", middleware.JWTAuth("testsecret"), handler.Profile)
	return router
}

func TestIntegration_RegisterAndLogin(t *testing.T) {
	// Set Gin to Test Mode.
	gin.SetMode(gin.TestMode)
	router := setupTestRouter()

	// 1. Test Registration
	registerPayload := handlers.RegisterInput{
		Username:    "inttestuser",
		Email:       "inttestuser@example.com",
		Password:    "inttestpassword",
		Preferences: map[string]interface{}{"diet": "vegan"},
	}
	body, err := json.Marshal(registerPayload)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, "Expected status code 201 on register")

	// TODO: Add tests for registration errors:
	// - Malformed JSON payload
	// - Missing required fields (e.g., missing password or email)
	// - Invalid email format

	// 2. Test Login
	loginPayload := map[string]string{
		"email":    "inttestuser@example.com",
		"password": "inttestpassword",
	}
	body, err = json.Marshal(loginPayload)
	assert.NoError(t, err)

	req = httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200 on login")

	var loginResp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NoError(t, err)
	token, ok := loginResp["token"]
	assert.True(t, ok, "Expected token in login response")
	assert.NotEmpty(t, token, "Token should not be empty")

	// TODO: Add tests for login errors:
	// - Incorrect password (already tested in TestIntegration_InvalidLogin below)
	// - Malformed JSON in login payload

	// 3. Test Profile (Protected Endpoint)
	req = httptest.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200 on profile endpoint")

	var profileResp models.User
	err = json.Unmarshal(w.Body.Bytes(), &profileResp)
	assert.NoError(t, err)
	assert.Equal(t, "inttestuser@example.com", profileResp.Email, "Profile email should match")

	// TODO: Add tests for profile endpoint errors:
	// - Missing or invalid token in Authorization header
	// - Expired or tampered token (simulate by altering the token)
}

func TestIntegration_InvalidLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupTestRouter()

	// First, register a user.
	registerPayload := handlers.RegisterInput{
		Username:    "inttestuser2",
		Email:       "inttestuser2@example.com",
		Password:    "validpassword",
		Preferences: map[string]interface{}{"diet": "vegetarian"},
	}
	body, err := json.Marshal(registerPayload)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Attempt to log in with the wrong password.
	loginPayload := map[string]string{
		"email":    "inttestuser2@example.com",
		"password": "wrongpassword",
	}
	body, err = json.Marshal(loginPayload)
	assert.NoError(t, err)

	req = httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status code 401 for invalid login")

	// TODO: Consider testing other invalid login scenarios:
	// - Non-existent user email
	// - Malformed login payload
}
