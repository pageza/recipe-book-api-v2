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
	log.Println("Initial database connection established")

	// In CI, skip migrations in the API container.
	log.Println("Skipping migrations in API container. Assuming migration container has applied schema.")

	// Poll until the "users" table exists.
	maxWait := 30 * time.Second
	interval := 2 * time.Second
	waited := time.Duration(0)
	for {
		if db.Migrator().HasTable(&models.User{}) {
			log.Println("users table detected.")
			break
		}
		tables, err := db.Migrator().GetTables()
		if err != nil {
			log.Printf("Failed to retrieve table list: %v", err)
		} else {
			log.Printf("Current tables in database: %v", tables)
		}
		if waited >= maxWait {
			log.Fatalf("users table does not exist after waiting %v", maxWait)
		}
		log.Printf("Waiting for users table to be created... waited %v", waited)
		time.Sleep(interval)
		waited += interval
	}
	log.Println("Database ready: users table exists")

	// Force a full reconnection by closing and resetting the connection pool.
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to retrieve underlying sql.DB: %v", err)
	}
	// Optionally, set max idle connections to 0 to flush the pool.
	sqlDB.SetMaxIdleConns(0)
	err = sqlDB.Close()
	if err != nil {
		log.Fatalf("failed to close the stale connection pool: %v", err)
	}
	log.Println("Stale DB connection pool closed. Reconnecting...")

	// Reconnect to the database.
	db, err = config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to reconnect to database: %v", err)
	}
	// Verify that the new connection sees the users table.
	if !db.Migrator().HasTable(&models.User{}) {
		log.Fatalf("reconnected DB does not see users table")
	}
	log.Println("New DB connection established and confirmed schema.")

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
