package users_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/handlers/users"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// dummyUserService implements service.UserService
type dummyUserService struct{}

func (d *dummyUserService) Register(user *models.User) error {
	// Simulate successful registration.
	return nil
}

func (d *dummyUserService) Login(email, password string) (*models.User, error) {
	// Simulate successful login: generate hash from the provided password and return a dummy user.
	hash, _ := utils.HashPassword(password)
	return &models.User{
		ID:           "dummy-id",
		Email:        email,
		PasswordHash: hash,
		Username:     "dummy",
		Preferences:  "{}",
	}, nil
}

func (d *dummyUserService) GetProfile(userID string) (*models.User, error) {
	return &models.User{
		ID:           userID,
		Email:        "dummy@example.com",
		Username:     "dummy",
		PasswordHash: "",
		Preferences:  "{}",
	}, nil
}

func TestRegisterAndLoginHandler(t *testing.T) {
	// Set Gin to test mode.
	gin.SetMode(gin.TestMode)

	// Use the dummy service.
	dummySvc := &dummyUserService{}
	handler := users.NewUserHandler(dummySvc, "testsecret")

	router := gin.Default()
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	// Test Registration
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

	// Test Login
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
