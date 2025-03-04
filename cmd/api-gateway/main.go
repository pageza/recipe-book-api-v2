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

	"go.uber.org/zap"

	_ "github.com/pageza/recipe-book-api-v2/docs"
	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	"github.com/pageza/recipe-book-api-v2/internal/handlers/recipes"
	"github.com/pageza/recipe-book-api-v2/internal/handlers/users"
	"github.com/pageza/recipe-book-api-v2/internal/middleware"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
	"github.com/pageza/recipe-book-api-v2/internal/routes"
	"github.com/pageza/recipe-book-api-v2/internal/service"
)

func main() {
	// Initialize the global logger for both main and middleware.
	if err := middleware.InitLogger(); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	// Make sure Zap's global logger is the same as the one in middleware.
	zap.ReplaceGlobals(middleware.Log)
	defer func() {
		if err := middleware.SyncLogger(); err != nil {
			log.Printf("failed to sync logger: %v", err)
		}
	}()

	// Load configuration.
	cfg, err := config.LoadConfig()
	if err != nil {
		zap.L().Fatal("failed to load config", zap.Error(err))
	}
	zap.L().Info("Configuration loaded",
		zap.String("PORT", cfg.Port),
		zap.String("DB_HOST", cfg.DBHost),
		zap.String("DB_NAME", cfg.DBName),
	)

	// Connect to the PostgreSQL database.
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		zap.L().Fatal("failed to connect to database", zap.Error(err))
	}

	// Retrieve underlying sql.DB to manage idle connections.
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

	// Auto-migrate the Recipe model to ensure the table exists.
	if err := db.AutoMigrate(&models.Recipe{}); err != nil {
		zap.L().Fatal("failed to auto-migrate Recipe model", zap.Error(err))
	}

	// Initialize repositories, services, and handlers for users.
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := users.NewUserHandler(userService, cfg.JWTSecret)

	// Initialize repositories, services, and handlers for recipes.
	recipeRepo := repository.NewRecipeRepository(db)
	recipeService := service.NewRecipeService(recipeRepo)
	recipeHandler := recipes.NewRecipeHandler(recipeService)

	// Combine handlers.
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
