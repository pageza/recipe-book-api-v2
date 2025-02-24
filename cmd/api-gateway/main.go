/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/
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

	// In CI, drop existing tables to ensure a clean state.
	// if os.Getenv("CI") == "true" {
	// 	log.Println("CI environment detected, dropping existing tables")
	// 	if err := db.Migrator().DropTable(&models.User{}, &models.Recipe{}); err != nil {
	// 		log.Fatalf("failed to drop tables: %v", err)
	// 	}
	// }

	// Run migrations for required models.
	err = db.AutoMigrate(&models.User{}, &models.Recipe{})
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Database migrations complete")
	time.Sleep(5 * time.Second) // give the DB time to finalize DDL changes

	if !db.Migrator().HasTable(&models.User{}) {
		log.Fatalf("users table does not exist after migration")
	}

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
