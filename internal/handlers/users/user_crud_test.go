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

// dummyUserService is a dummy implementation of the service.UserService interface
// for testing CRUD endpoints.
type dummyUserService struct{}

func (d *dummyUserService) Register(user *models.User) error {
	return nil
}

func (d *dummyUserService) Login(email, password string) (*models.User, error) {
	return &models.User{
		ID:       "dummy-id",
		Username: "dummyuser",
		Email:    email,
	}, nil
}

func (d *dummyUserService) GetProfile(userID string) (*models.User, error) {
	return &models.User{
		ID:       userID,
		Username: "dummyuser",
		Email:    "dummy@example.com",
	}, nil
}

func (d *dummyUserService) UpdateUser(user *models.User) error {
	// Simulate a successful update.
	return nil
}

func (d *dummyUserService) DeleteUser(userID string) error {
	// Simulate a successful deletion.
	return nil
}

func (d *dummyUserService) GetUserByEmail(email string) (*models.User, error) {
	return &models.User{
		ID:       "dummy-id",
		Username: "dummyuser",
		Email:    email,
	}, nil
}

// TestUpdateUserEndpoint tests the PUT /user endpoint.
func TestUpdateUserEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "testsecret"
	dummySvc := &dummyUserService{}
	userHandler := users.NewUserHandler(dummySvc, secret)

	router := gin.Default()
	// Assuming the update endpoint is registered at PUT /user.
	router.PUT("/user", userHandler.Update)

	// Prepare the update payload.
	updatePayload := map[string]string{
		"id":       "dummy-id",
		"username": "updateduser",
		"email":    "updated@example.com",
	}
	payloadBytes, err := json.Marshal(updatePayload)
	assert.NoError(t, err)

	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	// Generate a valid JWT token using the new extended signature.
	token, err := utils.GenerateJWT("dummy-id", "user", []string{}, secret)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check for a successful response.
	assert.Equal(t, http.StatusOK, w.Code)
	// You can add additional assertions based on your handler's response.
}

// TestDeleteUserEndpoint tests the DELETE /user/:id endpoint.
func TestDeleteUserEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "testsecret"
	dummySvc := &dummyUserService{}
	userHandler := users.NewUserHandler(dummySvc, secret)

	router := gin.Default()
	// Assuming the delete endpoint is registered at DELETE /user/:id.
	router.DELETE("/user/:id", userHandler.Delete)

	req, _ := http.NewRequest("DELETE", "/user/dummy-id", nil)
	token, err := utils.GenerateJWT("dummy-id", "user", []string{}, secret)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check that the deletion returns a successful status.
	assert.Equal(t, http.StatusOK, w.Code)
	// You can further validate the response if your endpoint returns a message or status.
}
