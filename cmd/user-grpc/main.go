package main

import (
	"net"

	grpcserver "github.com/pageza/recipe-book-api-v2/grpc/user"
	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	pb "github.com/pageza/recipe-book-api-v2/proto/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		zap.L().Fatal("Failed to load config", zap.Error(err))
	}

	// Initialize database connection
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		zap.L().Fatal("Failed to connect to the database", zap.Error(err))
	}

	// Initialize dependencies
	repo := repository.NewUserRepository(db) // ✅ Fixed missing *gorm.DB
	userSvc := service.NewUserService(repo)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, grpcserver.NewServer(userSvc))

	// Listen and serve
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		zap.L().Fatal("Failed to listen", zap.Error(err))
	}
	zap.L().Info("User gRPC Service running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		zap.L().Fatal("Failed to serve", zap.Error(err))
	}
}
