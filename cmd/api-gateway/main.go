/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

// cmd/api-gateway/main.go
package main

import (
	"log"
	"time"

	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	"github.com/pageza/recipe-book-api-v2/internal/handlers/recipes"
	"github.com/pageza/recipe-book-api-v2/internal/handlers/users"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
	"github.com/pageza/recipe-book-api-v2/internal/routes"
	"github.com/pageza/recipe-book-api-v2/internal/service"
)

func main() {
	// Load configuration and log it.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("Configuration loaded: PORT=%s, DB_HOST=%s, DB_NAME=%s", cfg.Port, cfg.DBHost, cfg.DBName)

	// Connect to the database.
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Database connection established")

	// Run migrations for required models.
	err = db.AutoMigrate(&models.User{}, &models.Recipe{})
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Database migrations initiated")

	// Wait until the "users" table is confirmed to exist.
	maxWait := 30 * time.Second
	interval := 2 * time.Second
	waited := time.Duration(0)
	for !db.Migrator().HasTable(&models.User{}) {
		if waited >= maxWait {
			log.Fatalf("users table does not exist after waiting %v", maxWait)
		}
		log.Println("Waiting for users table to be created...")
		time.Sleep(interval)
		waited += interval
	}
	log.Println("Database migrations complete: users table exists")

	// Initialize repositories, services, and handlers.
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := users.NewUserHandler(userService, cfg.JWTSecret)

	recipeRepo := repository.NewRecipeRepository(db)
	recipeService := service.NewRecipeService(recipeRepo)
	recipeHandler := recipes.NewRecipeHandler(recipeService)

	// Combine all handlers into a composite struct.
	h := &handlers.Handlers{
		User:   userHandler,
		Recipe: recipeHandler,
		// Add notifications handler when ready.
	}

	// Initialize the router using the composite handlers.
	r := routes.NewRouter(cfg, h)
	log.Println("Router initialized, ready to start server")

	// Start the server.
	addr := ":" + cfg.Port
	log.Printf("Starting API server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
