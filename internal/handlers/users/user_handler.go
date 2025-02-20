/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package users

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	service   service.UserService // or service.UserServiceInterface if that's your defined interface
	jwtSecret string
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(svc service.UserService, jwtSecret string) *UserHandler {
	return &UserHandler{
		service:   svc,
		jwtSecret: jwtSecret,
	}
}

type RegisterInput struct {
	Username    string                 `json:"username" binding:"required"`
	Email       string                 `json:"email" binding:"required,email"`
	Password    string                 `json:"password" binding:"required"`
	Preferences map[string]interface{} `json:"preferences"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password before storing it
	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	// Convert the map to a JSON string.
	prefBytes, err := json.Marshal(input.Preferences)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process preferences"})
		return
	}
	preferencesStr := string(prefBytes)

	// Create a new user model
	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashed, // store the hashed version
		Preferences:  preferencesStr,
	}

	// If no ID is set, generate one.
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	if err := h.service.Register(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := utils.GenerateJWT(user.ID, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	user, err := h.service.GetProfile(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
