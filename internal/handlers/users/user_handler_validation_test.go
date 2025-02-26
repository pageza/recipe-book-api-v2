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
	"github.com/pageza/recipe-book-api-v2/internal/service"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// errorUserService simulates error conditions for testing validation and error handling.
type errorUserService struct{}

func (e *errorUserService) Register(user *models.User) error {
	if user.Email == "" {
		return errors.Wrap(errors.New("email is required"), "registration error")
	}
	return errors.Wrap(errors.New("unexpected registration error"), "registration error")
}

func (e *errorUserService) Login(email, password string) (*models.User, error) {
	// Return a wrapped version of ErrInvalidCredentials.
	return nil, errors.Wrap(service.ErrInvalidCredentials, "login error")
}

func (e *errorUserService) GetProfile(userID string) (*models.User, error) {
	return nil, errors.Wrap(errors.New("user not found"), "profile error")
}

// Added stub for UpdateUser to satisfy the interface.
func (e *errorUserService) UpdateUser(user *models.User) error {
	return errors.Wrap(errors.New("update error"), "update error")
}

// Added stub for DeleteUser to satisfy the interface.
func (e *errorUserService) DeleteUser(userID string) error {
	return errors.Wrap(errors.New("delete error"), "delete error")
}

func (e *errorUserService) GetUserByEmail(email string) (*models.User, error) {
	return nil, errors.Wrap(errors.New("get user error: user not found"), "get user error")
}

// duplicateUserService simulates a duplicate registration scenario.
type duplicateUserService struct {
	registered bool
}

func (d *duplicateUserService) Register(user *models.User) error {
	if !d.registered {
		d.registered = true
		return nil
	}
	return errors.Wrap(service.ErrUserAlreadyExists, "duplicate registration error")
}

func (d *duplicateUserService) Login(email, password string) (*models.User, error) {
	return nil, errors.Wrap(errors.New("duplicate user login error"), "login error")
}

func (d *duplicateUserService) GetProfile(userID string) (*models.User, error) {
	return nil, errors.Wrap(errors.New("duplicate user profile error"), "profile error")
}

func (d *duplicateUserService) UpdateUser(user *models.User) error {
	return errors.Wrap(errors.New("duplicate update error"), "update error")
}

func (d *duplicateUserService) DeleteUser(userID string) error {
	return errors.Wrap(errors.New("duplicate delete error"), "delete error")
}

func (d *duplicateUserService) GetUserByEmail(email string) (*models.User, error) {
	return nil, errors.Wrap(errors.New("duplicate get user error"), "get user error")
}

// validUserService simulates a service that returns valid user data.
type validUserService struct{}

func (v *validUserService) Register(user *models.User) error {
	return nil
}

func (v *validUserService) Login(email, password string) (*models.User, error) {
	return &models.User{
		ID:       "valid-user-id",
		Username: "validuser",
		Email:    "valid@example.com",
	}, nil
}

func (v *validUserService) GetProfile(userID string) (*models.User, error) {
	return &models.User{
		ID:       userID,
		Username: "validuser",
		Email:    "valid@example.com",
	}, nil
}

func (v *validUserService) UpdateUser(user *models.User) error {
	return nil
}

func (v *validUserService) DeleteUser(userID string) error {
	return nil
}

func (v *validUserService) GetUserByEmail(email string) (*models.User, error) {
	return &models.User{
		ID:       "valid-id",
		Username: "validuser",
		Email:    "valid@example.com",
	}, nil
}

func TestRegisterValidation_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &errorUserService{}
	handler := users.NewUserHandler(svc, "testsecret")
	router := gin.Default()
	router.POST("/register", handler.Register)

	// Create payload with missing email.
	payload := users.RegisterInput{
		Username:    "testuser",
		Email:       "", // Missing email should trigger a validation error.
		Password:    "password",
		Preferences: map[string]interface{}{"diet": "vegan"},
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Expect a 400 Bad Request status due to missing email.
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status 400 for missing email")
}

func TestRegisterInvalidEmailFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &errorUserService{}
	handler := users.NewUserHandler(svc, "testsecret")
	router := gin.Default()
	router.POST("/register", handler.Register)

	// Create payload with an invalid email format.
	payload := users.RegisterInput{
		Username:    "testuser",
		Email:       "invalid-email", // Invalid format.
		Password:    "password",
		Preferences: map[string]interface{}{"diet": "vegan"},
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Gin binding should catch the invalid email format and return 400.
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status 400 for invalid email format")
}

func TestRegisterMissingPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &errorUserService{}
	handler := users.NewUserHandler(svc, "testsecret")
	router := gin.Default()
	router.POST("/register", handler.Register)

	// Create payload omitting the password field.
	payload := map[string]interface{}{
		"username":    "testuser",
		"email":       "testuser@example.com",
		"preferences": map[string]interface{}{"diet": "vegan"},
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Expect 400 due to missing password.
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status 400 for missing password")
}

func TestLoginErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &errorUserService{}
	handler := users.NewUserHandler(svc, "testsecret")
	router := gin.Default()
	router.POST("/login", handler.Login)

	// Prepare a login payload with invalid credentials.
	payload := map[string]string{
		"email":    "testuser@example.com",
		"password": "wrongpassword",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Expect a 401 Unauthorized status when credentials are invalid.
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status 401 for invalid login")
}

func TestLoginMissingPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &errorUserService{}
	handler := users.NewUserHandler(svc, "testsecret")
	router := gin.Default()
	router.POST("/login", handler.Login)

	// Create a login payload missing the password field.
	payload := map[string]string{
		"email": "testuser@example.com",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Expect a 400 status for missing password.
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status 400 for missing password in login")
}

func TestLoginInvalidEmailFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &errorUserService{}
	handler := users.NewUserHandler(svc, "testsecret")
	router := gin.Default()
	router.POST("/login", handler.Login)

	// Create a login payload with an invalid email format.
	payload := map[string]string{
		"email":    "invalid-email",
		"password": "password",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Expect a 400 status due to invalid email format.
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status 400 for invalid email format in login")
}

func TestRegisterDuplicate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &duplicateUserService{}
	handler := users.NewUserHandler(svc, "testsecret")
	router := gin.Default()
	router.POST("/register", handler.Register)

	// Create a valid registration payload.
	payload := users.RegisterInput{
		Username:    "dupuser",
		Email:       "dup@example.com",
		Password:    "password",
		Preferences: map[string]interface{}{"diet": "vegan"},
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	// First registration attempt should succeed with status 201.
	req1 := httptest.NewRequest("POST", "/register", bytes.NewReader(b))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code, "Expected first registration to succeed")

	// Second registration attempt should fail with a 409 Conflict.
	req2 := httptest.NewRequest("POST", "/register", bytes.NewReader(b))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code, "Expected duplicate registration to return conflict")
}

func TestGetProfileSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &validUserService{}
	handler := users.NewUserHandler(svc, "testsecret")
	router := gin.Default()

	// Inject claims as if added by JWT middleware.
	router.Use(func(c *gin.Context) {
		c.Set("userClaims", &utils.JWTClaims{UserID: "valid-user-id"})
		c.Next()
	})
	router.GET("/profile", handler.Profile)

	req := httptest.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 for successful profile fetch")

	// Updated structure reflecting flat JSON response:
	var resp struct {
		Status      string      `json:"status"`
		ID          string      `json:"id"`
		Email       string      `json:"email"`
		Username    string      `json:"username"`
		Preferences interface{} `json:"preferences"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "valid@example.com", resp.Email, "Profile email should match")
	assert.Equal(t, "validuser", resp.Username, "Profile username should match")
}

func TestRegisterMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &errorUserService{}
	handler := users.NewUserHandler(svc, "testsecret")
	router := gin.Default()
	router.POST("/register", handler.Register)

	// Send malformed JSON.
	malformedJSON := "{username: testuser, email: 'badjson'"
	req := httptest.NewRequest("POST", "/register", bytes.NewReader([]byte(malformedJSON)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Expect a 400 Bad Request status.
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status 400 for malformed JSON in registration")
}

func TestLoginMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &errorUserService{}
	handler := users.NewUserHandler(svc, "testsecret")
	router := gin.Default()
	router.POST("/login", handler.Login)

	// Send malformed JSON.
	malformedJSON := "{email: testuser@example.com, password: 'badjson'"
	req := httptest.NewRequest("POST", "/login", bytes.NewReader([]byte(malformedJSON)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Expect a 400 Bad Request status.
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status 400 for malformed JSON in login")
}
