/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service   service.UserService
	jwtSecret string
}

type RegisterInput struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	Preferences string `json:"preferences"`
}

func NewUserHandler(s service.UserService, jwtSecret string) *UserHandler {
	return &UserHandler{service: s, jwtSecret: jwtSecret}
}

func (h *UserHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("DEBUG: Register - failed to bind JSON: %v", err)
		return
	}
	log.Printf("DEBUG: Register - received input: %+v", input)

	// Hash the password before storing it
	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		log.Printf("DEBUG: Register - hashing failed: %v", err)
		return
	}
	log.Printf("DEBUG: Register - hashed password: %s", hashed)

	// Create a new user model
	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashed, // store the hashed version
		Preferences:  input.Preferences,
	}

	// If no ID is set, generate one.
	if user.ID == "" {
		user.ID = uuid.New().String()
		log.Printf("DEBUG: Register - generated new user ID: %s", user.ID)
	}

	if err := h.service.Register(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("DEBUG: Register - service error: %v", err)
		return
	}
	log.Printf("DEBUG: Register - user registered successfully for email: %s", input.Email)
	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("DEBUG: Login - failed to bind JSON: %v", err)
		return
	}
	log.Printf("DEBUG: Handler Login - received email: %s", input.Email)
	log.Printf("DEBUG: Handler Login - received password: %s", input.Password)

	user, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		log.Printf("DEBUG: Handler Login - login failed with error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := utils.GenerateJWT(user.ID, h.jwtSecret)
	if err != nil {
		log.Printf("DEBUG: Handler Login - failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	log.Printf("DEBUG: Handler Login - successfully generated token for email: %s", input.Email)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		log.Printf("DEBUG: Profile - no userID found in context")
		return
	}
	user, err := h.service.GetProfile(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		log.Printf("DEBUG: Profile - user not found for ID: %v", userID)
		return
	}
	log.Printf("DEBUG: Profile - returning user for ID: %v", userID)
	c.JSON(http.StatusOK, user)
}
