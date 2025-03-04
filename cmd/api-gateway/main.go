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
	"time"

	"go.uber.org/zap"

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
		zap.L().Fatal("Failed to load config", zap.Error(err))
	}
	zap.L().Info("Configuration loaded", zap.String("PORT", cfg.Port), zap.String("DB_HOST", cfg.DBHost), zap.String("DB_NAME", cfg.DBName))

	// Connect to the database.
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		zap.L().Fatal("failed to connect to database", zap.Error(err))
	}
	zap.L().Info("Initial database connection established")

	// Skip migrations in the API container.
	zap.L().Info("Skipping migrations in API container. Assuming migration container has applied schema.")

	// Wait (via raw SQL query) until the 'users' table is visible.
	maxWait := 30 * time.Second
	interval := 2 * time.Second
	waited := time.Duration(0)
	for {
		var count int64
		err := db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_name='users'").Scan(&count).Error
		if err != nil {
			zap.L().Warn("Failed to check table existence", zap.Error(err))
		} else {
			zap.L().Info("Raw check", zap.Int64("count", count))
			if count > 0 {
				break
			}
		}
		if waited >= maxWait {
			zap.L().Fatal("'users' table not found after waiting", zap.Duration("maxWait", maxWait))
		}
		zap.L().Info("Waiting for 'users' table...", zap.Duration("waited", waited))
		time.Sleep(interval)
		waited += interval
	}
	zap.L().Info("Database ready: 'users' table is visible.")

	// Force a full reconnection: close and reset the connection pool.
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Fatal("failed to retrieve underlying sql.DB", zap.Error(err))
	}
	// Flush idle connections.
	sqlDB.SetMaxIdleConns(0)
	if err := sqlDB.Close(); err != nil {
		zap.L().Fatal("failed to close connection pool", zap.Error(err))
	}
	zap.L().Info("Closed stale connection pool. Reconnecting...")

	// Reconnect to the database.
	db, err = config.ConnectDatabase(cfg)
	if err != nil {
		zap.L().Fatal("failed to reconnect to database", zap.Error(err))
	}
	// Verify that the new connection sees the 'users' table.
	var verifyCount int64
	if err := db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_name='users'").Scan(&verifyCount).Error; err != nil {
		zap.L().Fatal("failed to verify table existence", zap.Error(err))
	}
	if verifyCount == 0 {
		zap.L().Fatal("New DB connection does not see 'users' table")
	}
	zap.L().Info("New DB connection established and confirmed schema.")

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
	zap.L().Info("Router initialized, ready to start server")

	// Start the server.
	addr := ":" + cfg.Port
	zap.L().Info("Starting API server", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		zap.L().Fatal("failed to run server", zap.Error(err))
	}
}
