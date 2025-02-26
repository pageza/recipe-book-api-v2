/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package service

import (
	"errors"
	"log"

	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"
)

var (
	// ErrUserAlreadyExists is returned when a duplicate registration is attempted.
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidCredentials is returned when the login credentials are incorrect.
	ErrInvalidCredentials = errors.New("invalid credentials")
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

	if user.Email == "" {
		return errors.New("email cannot be empty")
	}

	if existing, _ := s.repo.GetUserByEmail(user.Email); existing != nil {
		log.Printf("Register: duplicate registration attempted for email: %s", user.Email)
		return ErrUserAlreadyExists
	}
	err := s.repo.CreateUser(user)
	if err != nil {
		log.Printf("Register: failed to create user (%s): %v", user.Email, err)
	} else {
		log.Printf("Register: user (%s) registered successfully", user.Email)
	}
	return err
}

// Login retrieves the user by email and verifies the password.
func (s *userService) Login(email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		// If the user is not found (or repository returns gorm.ErrRecordNotFound),
		// return the constant error.
		return nil, ErrUserNotFound
	}

	// Use the password hash check to compare the provided password.
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}
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
