/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

// @title Recipe Book API
// @version 1.0
// @description API documentation for the Recipe Book API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"log"
	"time"

	_ "github.com/pageza/recipe-book-api-v2/docs"
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
	log.Println("Initial database connection established")

	// Skip migrations in the API container.
	log.Println("Skipping migrations in API container. Assuming migration container has applied schema.")

	// Wait (via raw SQL query) until the 'users' table is visible.
	maxWait := 30 * time.Second
	interval := 2 * time.Second
	waited := time.Duration(0)
	for {
		var count int64
		err := db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_name='users'").Scan(&count).Error
		if err != nil {
			log.Printf("Failed to check table existence: %v", err)
		} else {
			log.Printf("Raw check: found %d 'users' table(s)", count)
			if count > 0 {
				break
			}
		}
		if waited >= maxWait {
			log.Fatalf("'users' table not found after waiting %v", maxWait)
		}
		log.Printf("Waiting for 'users' table... waited %v", waited)
		time.Sleep(interval)
		waited += interval
	}
	log.Println("Database ready: 'users' table is visible.")

	// Force a full reconnection: close and reset the connection pool.
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to retrieve underlying sql.DB: %v", err)
	}
	// Flush idle connections.
	sqlDB.SetMaxIdleConns(0)
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("failed to close connection pool: %v", err)
	}
	log.Println("Closed stale connection pool. Reconnecting...")

	// Reconnect to the database.
	db, err = config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to reconnect to database: %v", err)
	}
	// Verify that the new connection sees the 'users' table.
	var verifyCount int64
	if err := db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_name='users'").Scan(&verifyCount).Error; err != nil {
		log.Fatalf("failed to verify table existence: %v", err)
	}
	if verifyCount == 0 {
		log.Fatalf("New DB connection does not see 'users' table")
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
