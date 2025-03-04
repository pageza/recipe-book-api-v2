package main

import (
	"log"

	"go.uber.org/zap"

	grpcserver "github.com/pageza/recipe-book-api-v2/grpc"
	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
	"github.com/pageza/recipe-book-api-v2/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		zap.L().Fatal("Failed to load config", zap.Error(err))
	}

	// Connect to the database using config
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		zap.L().Fatal("Failed to connect to DB", zap.Error(err))
	}

	// Initialize repositories and services
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)

	recipeRepo := repository.NewRecipeRepository(db)
	recipeSvc := service.NewRecipeService(recipeRepo)

	notificationRepo := repository.NewNotificationRepository(db)
	storeEnabled := db != nil                                                         // ✅ Enable storage if DB is available
	notificationSvc := service.NewNotificationService(notificationRepo, storeEnabled) // ✅ Pass required args

	// Start the centralized gRPC server
	zap.L().Info("gRPC server started on port 50051")
	log.Fatal(grpcserver.StartGRPCServer(userSvc, recipeSvc, *notificationSvc)) // ✅ Pass value instead of pointer

}
