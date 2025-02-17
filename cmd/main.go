/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package main

import (
	"log"

	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	"github.com/pageza/recipe-book-api-v2/internal/middleware"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
	"github.com/pageza/recipe-book-api-v2/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration and connect to the database
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Initialize repositories, services, and handlers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, cfg.JWTSecret)

	// Initialize Gin router and apply middleware
	router := gin.Default()
	router.Use(middleware.Logger())

	// Public routes
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		protected.GET("/profile", userHandler.Profile)
	}

	// Start the server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
