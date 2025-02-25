package main

import (
	"log"

	grpcserver "github.com/pageza/recipe-book-api-v2/grpc"
	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
	"github.com/pageza/recipe-book-api-v2/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to the database using config
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
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
	log.Fatal(grpcserver.StartGRPCServer(userSvc, recipeSvc, *notificationSvc)) // ✅ Pass value instead of pointer

}
