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
	"net/http"
	"time"

	_ "github.com/pageza/recipe-book-api-v2/docs"
	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	"github.com/pageza/recipe-book-api-v2/internal/handlers/recipes"
	"github.com/pageza/recipe-book-api-v2/internal/handlers/users"
	"github.com/pageza/recipe-book-api-v2/internal/middleware"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
	"github.com/pageza/recipe-book-api-v2/internal/routes"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	"go.uber.org/zap"
)

func main() {
	// Initialize the zap logger
	if err := middleware.Init(); err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer middleware.Sync()

	middleware.Log.Info("Starting API Gateway")

	// Load configuration.
	cfg, err := config.LoadConfig()
	if err != nil {
		middleware.Log.Fatal("Failed to load config", zap.Error(err))
	}
	middleware.Log.Info("Configuration loaded", zap.String("PORT", cfg.Port), zap.String("DB_HOST", cfg.DBHost), zap.String("DB_NAME", cfg.DBName))

	// Connect to the database.
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		middleware.Log.Fatal("failed to connect to database", zap.Error(err))
	}
	middleware.Log.Info("Initial database connection established")

	// Skip migrations in the API container.
	middleware.Log.Info("Skipping migrations in API container. Assuming migration container has applied schema.")

	// Wait (via raw SQL query) until the 'users' table is visible.
	maxWait := 30 * time.Second
	interval := 2 * time.Second
	waited := time.Duration(0)
	for {
		var count int64
		err := db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_name='users'").Scan(&count).Error
		if err != nil {
			middleware.Log.Error("Failed to check table existence", zap.Error(err))
		} else {
			middleware.Log.Info("Raw check", zap.Int64("count", count))
			if count > 0 {
				break
			}
		}
		if waited >= maxWait {
			middleware.Log.Fatal("'users' table not found after waiting", zap.Duration("maxWait", maxWait))
		}
		middleware.Log.Info("Waiting for 'users' table...", zap.Duration("waited", waited))
		time.Sleep(interval)
		waited += interval
	}
	middleware.Log.Info("Database ready: 'users' table is visible.")

	// Force a full reconnection: close and reset the connection pool.
	sqlDB, err := db.DB()
	if err != nil {
		middleware.Log.Fatal("failed to retrieve underlying sql.DB", zap.Error(err))
	}
	// Flush idle connections.
	sqlDB.SetMaxIdleConns(0)
	if err := sqlDB.Close(); err != nil {
		middleware.Log.Fatal("failed to close connection pool", zap.Error(err))
	}
	middleware.Log.Info("Closed stale connection pool. Reconnecting...")

	// Reconnect to the database.
	db, err = config.ConnectDatabase(cfg)
	if err != nil {
		middleware.Log.Fatal("failed to reconnect to database", zap.Error(err))
	}
	// Verify that the new connection sees the 'users' table.
	var verifyCount int64
	if err := db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_name='users'").Scan(&verifyCount).Error; err != nil {
		middleware.Log.Fatal("failed to verify table existence", zap.Error(err))
	}
	if verifyCount == 0 {
		middleware.Log.Fatal("New DB connection does not see 'users' table")
	}
	middleware.Log.Info("New DB connection established and confirmed schema.")

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
	middleware.Log.Info("Router initialized, ready to start server")

	// Create HTTP server with graceful timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	middleware.Log.Info("Starting API server", zap.String("addr", ":"+cfg.Port))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		middleware.Log.Fatal("Server failed", zap.Error(err))
	}
}
