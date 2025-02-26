/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package service

import (
	"log"
	"net/http"

	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"
)

type AppError struct {
	Code int
	Msg  string
}

func (e *AppError) Error() string {
	return e.Msg
}

var (
	ErrUserAlreadyExists  = &AppError{Code: http.StatusConflict, Msg: "user already exists"}
	ErrUserNotFound       = &AppError{Code: http.StatusUnauthorized, Msg: "user not found"}
	ErrInvalidCredentials = &AppError{Code: http.StatusUnauthorized, Msg: "invalid credentials"}
	ErrEmailCannotBeEmpty = &AppError{Code: http.StatusBadRequest, Msg: "email cannot be empty"}
)

// UserService defines the business logic for user operations.
type UserService interface {
	Register(user *models.User) error
	Login(email, password string) (*models.User, error)
	GetProfile(userID string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(userID string) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new instance of the user service.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// Register creates a new user. It returns ErrUserAlreadyExists if the email is already registered.
func (s *userService) Register(user *models.User) error {
	log.Printf("[DEBUG][Service.Register] Registering user: %s", user.Email)

	if user.Email == "" {
		log.Printf("[DEBUG][Service.Register] Empty email encountered for user: %s", user.Email)
		return ErrEmailCannotBeEmpty
	}

	if existing, _ := s.repo.GetUserByEmail(user.Email); existing != nil {
		log.Printf("[DEBUG][Service.Register] Duplicate registration attempted for email: %s", user.Email)
		return ErrUserAlreadyExists
	}

	err := s.repo.CreateUser(user)
	if err != nil {
		log.Printf("[DEBUG][Service.Register] Failed to create user (%s): %v", user.Email, err)
	} else {
		log.Printf("[DEBUG][Service.Register] User (%s) registered successfully", user.Email)
	}
	return err
}

// Login retrieves the user by email and verifies the password.
func (s *userService) Login(email, password string) (*models.User, error) {
	log.Printf("[DEBUG][Service.Login] Attempting to find user by email: %s", email)
	user, err := s.repo.GetUserByEmail(email)
	if err != nil || user == nil {
		log.Printf("[DEBUG][Service.Login] User not found for email: %s, error: %v", email, err)
		return nil, ErrUserNotFound // now a generic 401 error
	}
	log.Printf("[DEBUG][Service.Login] User found for email: %s", email)

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		log.Printf("[DEBUG][Service.Login] Password mismatch for user %s", email)
		return nil, ErrInvalidCredentials
	}
	log.Printf("[DEBUG][Service.Login] Password verified for user %s", email)
	return user, nil
}

// GetProfile retrieves a user's profile by ID. If the user is not found, it returns ErrUserNotFound.
func (s *userService) GetProfile(userID string) (*models.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		log.Printf("GetProfile: user with ID (%s) not found: %v", userID, err)
		return nil, ErrUserNotFound
	}
	log.Printf("GetProfile: retrieved profile for user ID: %s", userID)
	return user, nil
}

// GetUserByEmail retrieves a user's profile by email. If the user is not found, it returns ErrUserNotFound.
func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		log.Printf("GetUserByEmail: user with email (%s) not found: %v", email, err)
		return nil, ErrUserNotFound
	}
	log.Printf("GetUserByEmail: retrieved profile for user email: %s", email)
	return user, nil
}

// UpdateUser updates an existing user. If a new password is provided, it will be hashed.
func (s *userService) UpdateUser(updated *models.User) error {
	// Confirm the user exists
	existing, err := s.repo.GetUserByID(updated.ID)
	if err != nil {
		return ErrUserNotFound
	}

	// Re-hash the new password if provided, otherwise keep the old hash.
	if updated.PasswordHash != "" {
		newHash, err := utils.HashPassword(updated.PasswordHash)
		if err != nil {
			log.Printf("UpdateUser: could not hash new password: %v", err)
			return err
		}
		updated.PasswordHash = newHash
	} else {
		updated.PasswordHash = existing.PasswordHash
	}

	// If email, username, or preferences are empty, keep the existing values.
	if updated.Email == "" {
		updated.Email = existing.Email
	}
	if updated.Username == "" {
		updated.Username = existing.Username
	}
	if updated.Preferences == "" {
		updated.Preferences = existing.Preferences
	}

	err = s.repo.UpdateUser(updated)
	if err != nil {
		log.Printf("UpdateUser: failed to update user %s: %v", updated.ID, err)
		return err
	}
	log.Printf("UpdateUser: user %s updated successfully", updated.ID)
	return nil
}

// DeleteUser removes a user by their ID.
func (s *userService) DeleteUser(userID string) error {
	// Ensure the user exists before deletion.
	_, err := s.repo.GetUserByID(userID)
	if err != nil {
		return ErrUserNotFound
	}
	err = s.repo.DeleteUser(userID)
	if err != nil {
		log.Printf("DeleteUser: failed to delete user %s: %v", userID, err)
		return err
	}
	log.Printf("DeleteUser: user %s deleted successfully", userID)
	return nil
}
