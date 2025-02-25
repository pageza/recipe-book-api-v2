/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/
// cmd/api-gateway/main.go
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

	// Set the search_path explicitly on this connection.
	if err := db.Exec("SET search_path TO public").Error; err != nil {
		log.Fatalf("failed to set search_path: %v", err)
	}

	// Wait (poll) until the "users" table is visible via a raw query.
	maxWait := 30 * time.Second
	interval := 2 * time.Second
	waited := time.Duration(0)
	for {
		var count int64
		err = db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_name='users'").Scan(&count).Error
		if err != nil {
			log.Printf("Failed to check table existence: %v", err)
		} else {
			log.Printf("Raw check: found %d 'users' table(s)", count)
			if count > 0 {
				break
			}
		}
		if waited >= maxWait {
			log.Fatalf("users table not visible after waiting %v", maxWait)
		}
		log.Printf("Waiting for users table to be visible... waited %v", waited)
		time.Sleep(interval)
		waited += interval
	}
	log.Println("Database ready: 'users' table is visible.")

	// (Optional) Force a reconnection: you might want to fully reset the pool.
	// For example, if needed, close the underlying sql.DB and reconnect.
	// [This code is commented out; uncomment if you suspect connection pooling issues.]
	/*
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("failed to get underlying sql.DB: %v", err)
		}
		sqlDB.Close()
		log.Println("Closed old connection pool; reconnecting...")
		db, err = config.ConnectDatabase(cfg)
		if err != nil {
			log.Fatalf("failed to reconnect: %v", err)
		}
		if err := db.Exec("SET search_path TO public").Error; err != nil {
			log.Fatalf("failed to set search_path on new connection: %v", err)
		}
	*/

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
