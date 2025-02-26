/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"
	"github.com/pkg/errors"
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

// Register creates a new user. It returns an appropriate wrapped error if the email is empty or already exists.
func (s *userService) Register(user *models.User) error {
	log.Printf("[DEBUG][Service.Register] Registering user: %s", user.Email)

	if user.Email == "" {
		log.Printf("[DEBUG][Service.Register] Empty email encountered for user: %s", user.Email)
		return errors.Wrap(ErrEmailCannotBeEmpty, "registration failed: empty email")
	}

	if existing, _ := s.repo.GetUserByEmail(user.Email); existing != nil {
		log.Printf("[DEBUG][Service.Register] Duplicate registration attempted for email: %s", user.Email)
		return errors.Wrap(ErrUserAlreadyExists, "registration failed: duplicate email")
	}

	err := s.repo.CreateUser(user)
	if err != nil {
		log.Printf("[DEBUG][Service.Register] Failed to create user (%s): %v", user.Email, err)
		return errors.Wrap(err, "registration failed: create user error")
	}

	log.Printf("[DEBUG][Service.Register] User (%s) registered successfully", user.Email)
	return nil
}

// Login retrieves the user by email and verifies the password.
func (s *userService) Login(email, password string) (*models.User, error) {
	log.Printf("[DEBUG][Service.Login] Attempting to find user by email: %s", email)
	user, err := s.repo.GetUserByEmail(email)
	if err != nil || user == nil {
		log.Printf("[DEBUG][Service.Login] User not found for email: %s, error: %v", email, err)
		return nil, errors.Wrap(ErrUserNotFound, fmt.Sprintf("login failed: user not found for email %s", email))
	}
	log.Printf("[DEBUG][Service.Login] User found for email: %s", email)

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		log.Printf("[DEBUG][Service.Login] Password mismatch for user %s", email)
		return nil, errors.Wrap(ErrInvalidCredentials, "login failed: invalid credentials")
	}
	log.Printf("[DEBUG][Service.Login] Password verified for user %s", email)
	return user, nil
}

// GetProfile retrieves a user's profile by ID.
func (s *userService) GetProfile(userID string) (*models.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil || user == nil {
		log.Printf("GetProfile: user with ID (%s) not found: %v", userID, err)
		return nil, errors.Wrap(ErrUserNotFound, fmt.Sprintf("profile retrieval failed for userID %s", userID))
	}
	log.Printf("GetProfile: retrieved profile for user ID: %s", userID)
	return user, nil
}

// GetUserByEmail retrieves a user's profile by email.
func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil || user == nil {
		log.Printf("GetUserByEmail: user with email (%s) not found: %v", email, err)
		return nil, errors.Wrap(ErrUserNotFound, fmt.Sprintf("failed to get user by email %s", email))
	}
	log.Printf("GetUserByEmail: retrieved profile for user email: %s", email)
	return user, nil
}

// UpdateUser updates an existing user.
func (s *userService) UpdateUser(updated *models.User) error {
	existing, err := s.repo.GetUserByID(updated.ID)
	if err != nil || existing == nil {
		return errors.Wrap(ErrUserNotFound, fmt.Sprintf("update failed: user %s not found", updated.ID))
	}

	if updated.PasswordHash != "" {
		newHash, err := utils.HashPassword(updated.PasswordHash)
		if err != nil {
			log.Printf("UpdateUser: could not hash new password: %v", err)
			return errors.Wrap(err, "update failed: password hash error")
		}
		updated.PasswordHash = newHash
	} else {
		updated.PasswordHash = existing.PasswordHash
	}

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
		return errors.Wrap(err, fmt.Sprintf("update failed for user %s", updated.ID))
	}
	log.Printf("UpdateUser: user %s updated successfully", updated.ID)
	return nil
}

// DeleteUser removes a user by their ID.
func (s *userService) DeleteUser(userID string) error {
	_, err := s.repo.GetUserByID(userID)
	if err != nil {
		return errors.Wrap(ErrUserNotFound, fmt.Sprintf("delete failed: user %s not found", userID))
	}
	err = s.repo.DeleteUser(userID)
	if err != nil {
		log.Printf("DeleteUser: failed to delete user %s: %v", userID, err)
		return errors.Wrap(err, fmt.Sprintf("delete failed for user %s", userID))
	}
	log.Printf("DeleteUser: user %s deleted successfully", userID)
	return nil
}
