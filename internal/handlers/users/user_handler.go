/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	service   service.UserService
	jwtSecret string
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(svc service.UserService, jwtSecret string) *UserHandler {
	return &UserHandler{
		service:   svc,
		jwtSecret: jwtSecret,
	}
}

// RegisterInput is the input payload for registering a new user.
// swagger:model RegisterInput
type RegisterInput struct {
	Username    string                 `json:"username" binding:"required"`
	Email       string                 `json:"email" binding:"required,email"`
	Password    string                 `json:"password" binding:"required"`
	Preferences map[string]interface{} `json:"preferences"`
}

/*
@Summary Register a new user
@Description Registers a new user in the system.
@Tags Users
@Accept  json
@Produce json
@Param register body RegisterInput true "User registration data"
@Success 201 {object} map[string]string "User registered successfully"
@Failure 400 {object} map[string]string "Bad request"
@Router /register [post]
*/
func (h *UserHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that email is provided.
	if input.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email cannot be empty"})
		return
	}

	// Hash the password before storing it.
	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	// Convert preferences to a JSON string.
	prefBytes, err := json.Marshal(input.Preferences)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process preferences"})
		return
	}
	preferencesStr := string(prefBytes)

	// Create a new user model.
	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashed,
		Preferences:  preferencesStr,
	}

	// Generate an ID if none is set.
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Attempt to register the user via the service layer.
	if err := h.service.Register(user); err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

/*
@Summary Login
@Description Authenticates a user with email and password and returns a JWT token.
@Tags Users
@Accept  json
@Produce json
@Param login body object{email=string,password=string} true "User login credentials"
@Success 200 {object} map[string]string "JWT token"
@Failure 400 {object} map[string]string "Bad request"
@Failure 401 {object} map[string]string "Unauthorized"
@Router /login [post]
*/
func (h *UserHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		// Check for invalid credentials.
		if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	token, err := utils.GenerateJWT(user.ID, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

/*
@Summary Get Profile
@Description Retrieves the profile of the authenticated user.
@Tags Users
@Produce json
@Success 200 {object} models.User "User profile information"
@Failure 401 {object} map[string]string "Unauthorized"
@Router /profile [get]
*/
func (h *UserHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	user, err := h.service.GetProfile(userID.(string))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}
