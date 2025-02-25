package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	userhandler "github.com/pageza/recipe-book-api-v2/internal/handlers/users"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/routes/protectedroutes"
	"github.com/pageza/recipe-book-api-v2/internal/routes/publicroutes"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// dummyService is a minimal implementation of service.UserService for router testing.
type dummyService struct{}

func (d *dummyService) Register(user *models.User) error { return nil }
func (d *dummyService) Login(email, password string) (*models.User, error) {
	return &models.User{
		ID:           "dummy-id",
		Email:        email,
		Username:     "dummyuser",
		PasswordHash: "dummy-hash",
		Preferences:  "{}",
	}, nil
}
func (d *dummyService) GetProfile(userID string) (*models.User, error) {
	return &models.User{
		ID:          userID,
		Email:       "dummy@example.com",
		Username:    "dummyuser",
		Preferences: "{}",
	}, nil
}

// New stub implementations to satisfy the interface:
func (d *dummyService) UpdateUser(user *models.User) error {
	// Simulate a successful update.
	return nil
}

func (d *dummyService) DeleteUser(userID string) error {
	// Simulate a successful deletion.
	return nil
}

func (d *dummyService) GetUserByEmail(email string) (*models.User, error) {
	// Return a dummy user based on the provided email.
	return &models.User{
		ID:       "dummy-id",
		Username: "dummyuser",
		Email:    email,
	}, nil
}

// newDummyHandlers returns a composite handlers.Handlers with a real user handler
// constructed using the dummyService.
func newDummyHandlers() *handlers.Handlers {
	// Use the real constructor from the userhandler package.
	uh := userhandler.NewUserHandler(&dummyService{}, "testsecret")
	return &handlers.Handlers{
		User: uh,
	}
}

// setupRouter creates a Gin router and registers both public and protected routes.
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Register public routes using the publicroutes package.
	publicroutes.Register(router, newDummyHandlers())

	// Create a dummy configuration with a JWT secret.
	cfg := &config.Config{
		JWTSecret: "testsecret",
	}
	// Register protected routes using the protectedroutes package.
	protectedroutes.Register(router, cfg, newDummyHandlers())

	return router
}

func TestPublicRoutes(t *testing.T) {
	router := setupRouter()

	// Test the /register endpoint.
	regPayload := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "pass",
	}
	regBytes, err := json.Marshal(regPayload)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(regBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// Expect a 201 Created response.
	assert.Equal(t, http.StatusCreated, w.Code)

	var regResp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &regResp)
	assert.NoError(t, err)
	// The real handler returns "user registered".
	assert.Equal(t, "user registered", regResp["message"])

	// Test the /login endpoint.
	loginPayload := map[string]string{
		"email":    "test@example.com",
		"password": "pass",
	}
	loginBytes, err := json.Marshal(loginPayload)
	assert.NoError(t, err)

	req = httptest.NewRequest("POST", "/login", bytes.NewReader(loginBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// Expect a 200 OK response.
	assert.Equal(t, http.StatusOK, w.Code)

	var loginResp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NoError(t, err)
	// Expect the token to match what is generated via utils.GenerateJWT.
	expectedToken, err := utils.GenerateJWT("dummy-id", "testsecret")
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, loginResp["token"])
}

func TestProtectedRoutes(t *testing.T) {
	router := setupRouter()

	// Attempt to access the protected /profile endpoint without a token.
	req := httptest.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Generate a valid token using utils.GenerateJWT.
	token, err := utils.GenerateJWT("dummy-id", "testsecret")
	assert.NoError(t, err)

	// Access the protected /profile endpoint with a valid token.
	req = httptest.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal the returned user object.
	var userResp models.User
	err = json.Unmarshal(w.Body.Bytes(), &userResp)
	assert.NoError(t, err)
	assert.Equal(t, "dummy-id", userResp.ID)
	assert.Equal(t, "dummy@example.com", userResp.Email)
	assert.Equal(t, "dummyuser", userResp.Username)
	assert.Equal(t, "{}", userResp.Preferences)
}
