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

// Login verifies user credentials. It returns ErrUserNotFound if the user does not exist,
// and ErrInvalidCredentials if the password does not match.
func (s *userService) Login(email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		log.Printf("Login: user with email (%s) not found: %v", email, err)
		return nil, ErrUserNotFound
	}
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		log.Printf("Login: invalid credentials for email: %s", email)
		return nil, ErrInvalidCredentials
	}
	log.Printf("Login: user (%s) logged in successfully", email)
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
