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

type UserService interface {
	Register(user *models.User) error
	Login(email, password string) (*models.User, error)
	GetProfile(userID string) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(user *models.User) error {
	log.Printf("DEBUG: Service Register - checking if user exists for email: %s", user.Email)
	if existing, _ := s.repo.GetUserByEmail(user.Email); existing != nil {
		log.Printf("DEBUG: Service Register - user already exists for email: %s", user.Email)
		return errors.New("user already exists")
	}
	err := s.repo.CreateUser(user)
	if err != nil {
		log.Printf("DEBUG: Service Register - error creating user: %v", err)
	} else {
		log.Printf("DEBUG: Service Register - user created successfully for email: %s", user.Email)
	}
	return err
}

func (s *userService) Login(email, password string) (*models.User, error) {
	log.Printf("DEBUG: Service Login - attempting login for email: %s", email)
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		log.Printf("DEBUG: Service Login - user not found for email: %s, error: %v", email, err)
		return nil, err
	}
	log.Printf("DEBUG: Service Login - incoming password: %s", password)
	log.Printf("DEBUG: Service Login - stored hash: %s", user.PasswordHash)
	match := utils.CheckPasswordHash(password, user.PasswordHash)
	log.Printf("DEBUG: Service Login - password match result: %v", match)
	if !match {
		log.Printf("DEBUG: Service Login - password check failed for email: %s", email)
		return nil, errors.New("invalid credentials")
	}
	log.Printf("DEBUG: Service Login - password check succeeded for email: %s", email)
	return user, nil
}

func (s *userService) GetProfile(userID string) (*models.User, error) {
	log.Printf("DEBUG: Service GetProfile - retrieving user for ID: %s", userID)
	return s.repo.GetUserByID(userID)
}
