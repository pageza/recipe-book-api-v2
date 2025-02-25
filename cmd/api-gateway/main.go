/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/
package main

import (
	"log"

	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	"github.com/pageza/recipe-book-api-v2/internal/handlers/recipes"
	"github.com/pageza/recipe-book-api-v2/internal/handlers/users"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
	"github.com/pageza/recipe-book-api-v2/internal/routes"
	"github.com/pageza/recipe-book-api-v2/internal/service"
)

func main() {
	// Load configuration.
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

	// No migration logic is performed here.
	log.Println("Assuming database schema has been set by migration container.")

	// Initialize repositories, services, and handlers.
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := users.NewUserHandler(userService, cfg.JWTSecret)

	recipeRepo := repository.NewRecipeRepository(db)
	recipeService := service.NewRecipeService(recipeRepo)
	recipeHandler := recipes.NewRecipeHandler(recipeService)

	h := &handlers.Handlers{
		User:   userHandler,
		Recipe: recipeHandler,
	}

	// Initialize the router.
	r := routes.NewRouter(cfg, h)
	log.Println("Router initialized, ready to start server")

	// Start the server.
	addr := ":" + cfg.Port
	log.Printf("Starting API server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
