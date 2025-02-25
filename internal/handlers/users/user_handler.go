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

	token, err := utils.GenerateJWT(user.ID, "user", []string{"read:profile"}, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Profile returns the profile of the authenticated user with additional information.
func (h *UserHandler) Profile(c *gin.Context) {
	// Get the extended claims from context.
	value, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	claims, ok := value.(*utils.JWTClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	user, err := h.service.GetProfile(claims.UserID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	// Return a structured JSON response, including extended claims.
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user": map[string]interface{}{
			"id":          user.ID,
			"username":    user.Username,
			"email":       user.Email,
			"role":        claims.Role,
			"permissions": claims.Permissions,
		},
	})
}

// RequestPasswordReset handles generating a password reset token for a user.
func (h *UserHandler) RequestPasswordReset(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, email required"})
		return
	}

	// Dummy token generation (in production, generate a secure token and store it)
	resetToken := uuid.New().String()

	// TODO: Store the reset token with expiration and send it via email.

	// For now, we return the token (for testing purposes)
	c.JSON(http.StatusOK, gin.H{"resetToken": resetToken})
}

// ResetPassword handles the password update using the reset token.
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Email       string `json:"email"`
		ResetToken  string `json:"resetToken"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.ResetToken == "" || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, missing fields"})
		return
	}

	// Dummy validation: in production, validate the resetToken matches what was stored.
	if req.ResetToken != "expected-dummy-token" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid reset token"})
		return
	}

	// Look up the user by email.
	user, err := h.service.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Set the new password. In a real app, you’d also re-hash the password.
	user.PasswordHash = req.NewPassword
	if err := h.service.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}
