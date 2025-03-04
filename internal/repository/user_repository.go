/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package repository

import (
	"go.uber.org/zap"

	"github.com/pageza/recipe-book-api-v2/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	// Migrate the schema; in production, use migration scripts.
	if err := db.AutoMigrate(&models.User{}); err != nil {
		zap.L().Warn("AutoMigrate error", zap.Error(err))
	}
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	err := r.db.Create(user).Error
	if err != nil {
	} else {
	}
	return err
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
