package main

import (
	"log"

	"net"

	grpcserver "github.com/pageza/recipe-book-api-v2/grpc/recipe"
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

	// Initialize dependencies
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	repo := repository.NewRecipeRepository(db) // ✅ Pass the actual DB instance

	recipeSvc := service.NewRecipeService(repo)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterRecipeServiceServer(grpcServer, grpcserver.NewServer(recipeSvc))

	// Listen and serve
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Recipe gRPC Service running on port 50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
