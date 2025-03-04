/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package service

import (
	"errors"

	"go.uber.org/zap"

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

	if user.Email == "" {
		return errors.New("email cannot be empty")
	}

	if existing, _ := s.repo.GetUserByEmail(user.Email); existing != nil {
		zap.L().Warn("Register: duplicate registration attempted", zap.String("email", user.Email))
		return ErrUserAlreadyExists
	}
	err := s.repo.CreateUser(user)
	if err != nil {
		zap.L().Error("Register: failed to create user", zap.String("email", user.Email), zap.Error(err))
	} else {
		zap.L().Info("Register: user registered successfully", zap.String("email", user.Email))
	}
	return err
}

// Login verifies user credentials. It returns ErrUserNotFound if the user does not exist,
// and ErrInvalidCredentials if the password does not match.
func (s *userService) Login(email, password string) (*models.User, error) {

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		zap.L().Warn("Login: user with email not found", zap.String("email", email), zap.Error(err))
		return nil, ErrUserNotFound
	}
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		zap.L().Warn("Login: invalid credentials for email", zap.String("email", email))
		return nil, ErrInvalidCredentials
	}
	zap.L().Info("Login: user logged in successfully", zap.String("email", email))
	return user, nil
}

// GetProfile retrieves a user's profile by ID. If the user is not found, it returns ErrUserNotFound.
func (s *userService) GetProfile(userID string) (*models.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		zap.L().Warn("GetProfile: user with ID not found", zap.String("userID", userID), zap.Error(err))
		return nil, ErrUserNotFound
	}
	zap.L().Info("GetProfile: retrieved profile for user ID", zap.String("userID", userID))
	return user, nil
}
