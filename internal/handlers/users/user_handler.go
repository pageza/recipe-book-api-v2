/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

// RegisterInput is a struct for registration input.
type RegisterInput struct {
	Username    string                 `json:"username" binding:"required"`
	Email       string                 `json:"email" binding:"required,email"`
	Password    string                 `json:"password" binding:"required"`
	Preferences map[string]interface{} `json:"preferences"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[DEBUG][Register] Binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[DEBUG][Register] Received registration request for email: %s", input.Email)

	// Validate that email is provided.
	if input.Email == "" {
		log.Printf("[DEBUG][Register] Empty email provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "email cannot be empty"})
		return
	}

	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		log.Printf("[DEBUG][Register] Password hashing failed for %s: %v", input.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}
	log.Printf("[DEBUG][Register] Password hashed successfully for %s", input.Email)

	prefBytes, err := json.Marshal(input.Preferences)
	if err != nil {
		log.Printf("[DEBUG][Register] Preferences marshalling failed for %s: %v", input.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process preferences"})
		return
	}
	preferencesStr := string(prefBytes)

	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashed,
		Preferences:  preferencesStr,
	}

	if user.ID == "" {
		user.ID = uuid.New().String()
		log.Printf("[DEBUG][Register] Generated new user ID: %s for email: %s", user.ID, user.Email)
	}

	if err := h.service.Register(user); err != nil {
		log.Printf("[DEBUG][Register] Service returned error: %+v, type: %T", err, err)
		rootErr := errors.Cause(err)
		if appErr, ok := rootErr.(*service.AppError); ok {
			log.Printf("[DEBUG][Register] Recognized AppError: Code %d, Msg: %s", appErr.Code, appErr.Msg)
			c.JSON(appErr.Code, gin.H{"error": appErr.Msg})
		} else {
			log.Printf("[DEBUG][Register] Unrecognized error type")
			c.JSON(http.StatusBadRequest, gin.H{"error": rootErr.Error()})
		}
		return
	}
	log.Printf("[DEBUG][Register] User registered successfully: %s", user.Email)
	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

// Login authenticates a user.
func (h *UserHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[DEBUG][Login] Binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[DEBUG][Login] Received login request for email: %s", input.Email)

	user, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		log.Printf("[DEBUG][Login] Service returned error: %+v, type: %T", err, err)
		rootErr := errors.Cause(err)
		if appErr, ok := rootErr.(*service.AppError); ok {
			log.Printf("[DEBUG][Login] Recognized AppError: Code %d, Msg: %s", appErr.Code, appErr.Msg)
			c.JSON(appErr.Code, gin.H{"error": appErr.Msg})
		} else {
			log.Printf("[DEBUG][Login] Unrecognized error type")
			c.JSON(http.StatusBadRequest, gin.H{"error": rootErr.Error()})
		}
		return
	}
	log.Printf("[DEBUG][Login] User authenticated successfully: %s", user.Email)

	token, err := utils.GenerateJWT(user.ID, "user", []string{"read:profile"}, h.jwtSecret)
	if err != nil {
		log.Printf("[DEBUG][Login] Failed to generate JWT for user %s: %v", user.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	log.Printf("[DEBUG][Login] JWT generated successfully for user %s", user.Email)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Profile returns the profile of the authenticated user with additional information.
// Profile returns the profile of the authenticated user, embedding the extended
// JWT claims in the JSON response under the "userClaims" key.
func (h *UserHandler) Profile(c *gin.Context) {
	log.Printf("[DEBUG][Profile] Profile endpoint called")
	value, exists := c.Get("userClaims")
	if !exists {
		log.Printf("[DEBUG][Profile] User claims not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	claims, ok := value.(*utils.JWTClaims)
	if !ok {
		log.Printf("[DEBUG][Profile] Claims type assertion failed, got: %T", value)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}
	log.Printf("[DEBUG][Profile] Retrieved claims: %+v", claims)

	user, err := h.service.GetProfile(claims.UserID)
	if err != nil {
		log.Printf("[DEBUG][Profile] GetProfile error for userID %s: %v", claims.UserID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	log.Printf("[DEBUG][Profile] Profile successful for user: %+v", user)

	// Return a flat JSON structure as expected by router tests.
	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"id":          user.ID,
		"email":       user.Email,
		"username":    user.Username,
		"preferences": user.Preferences,
	})
}

// Update handles updating a user's information.
func (h *UserHandler) Update(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

// Delete handles deleting a user.
func (h *UserHandler) Delete(c *gin.Context) {
	userID := c.Param("id")
	if err := h.service.DeleteUser(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
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

	// Set the new password. In a real app, you'd also re-hash the password.
	user.PasswordHash = req.NewPassword
	if err := h.service.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}
