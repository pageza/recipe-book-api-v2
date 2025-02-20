/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package service

import (
	"errors"

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
	if existing, _ := s.repo.GetUserByEmail(user.Email); existing != nil {
		return errors.New("user already exists")
	}
	err := s.repo.CreateUser(user)
	if err != nil {
	} else {
	}
	return err
}

func (s *userService) Login(email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	match := utils.CheckPasswordHash(password, user.PasswordHash)
	if !match {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *userService) GetProfile(userID string) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}
