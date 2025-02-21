package main

import (
	"log"
	"net"

	grpcserver "github.com/pageza/recipe-book-api-v2/grpc/notification"
	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	pb "github.com/pageza/recipe-book-api-v2/proto/proto"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Database connection (optional)
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Println("Warning: No database connection. Running in log-only mode.")
		db = nil // ✅ This prevents DB errors in log-only mode
	}

	// Initialize repository & service
	repo := repository.NewNotificationRepository(db)
	storeEnabled := db != nil // ✅ Enable storage if DB is connected
	notificationSvc := service.NewNotificationService(repo, storeEnabled)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, grpcserver.NewServer(notificationSvc))

	// Listen and serve
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Notification gRPC Service running on port 50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
