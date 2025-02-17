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
	log.Printf("DEBUG: Creating user with email: %s", user.Email)
	err := r.db.Create(user).Error
	if err != nil {
		log.Printf("DEBUG: Error creating user: %v", err)
	} else {
		log.Printf("DEBUG: User created successfully with email: %s", user.Email)
	}
	return err
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	log.Printf("DEBUG: Querying user by email: %s", email)
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Printf("DEBUG: GetUserByEmail error for %s: %v", email, err)
		return nil, err
	}
	log.Printf("DEBUG: Found user for email: %s with hash: %s", email, user.PasswordHash)
	return &user, nil
}

func (r *userRepository) GetUserByID(id string) (*models.User, error) {
	log.Printf("DEBUG: Querying user by ID: %s", id)
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		log.Printf("DEBUG: GetUserByID error for %s: %v", id, err)
		return nil, err
	}
	log.Printf("DEBUG: Found user for ID: %s", id)
	return &user, nil
}
