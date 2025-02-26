/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package repository

import (
	"log"

	"github.com/pageza/recipe-book-api-v2/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(userID string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	// Migrate the schema; in production, use migration scripts.
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Printf("DEBUG: AutoMigrate error: %v", err)
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

func (r *userRepository) UpdateUser(user *models.User) error {
	err := r.db.Save(user).Error
	if err != nil {
		log.Printf("UpdateUser: failed to update user %s, error: %v", user.ID, err)
		return err
	}
	return nil
}

func (r *userRepository) DeleteUser(userID string) error {
	result := r.db.Delete(&models.User{}, userID)
	if result.Error != nil {
		log.Printf("DeleteUser: failed delete for user %s, error: %v", userID, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Printf("DeleteUser: no user found with ID %s", userID)
		return gorm.ErrRecordNotFound
	}
	return nil
}
