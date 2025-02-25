package grpcserver

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/pageza/recipe-book-api-v2/proto/proto" // ✅ Unified import

	grpcNotification "github.com/pageza/recipe-book-api-v2/grpc/notification"
	grpcRecipe "github.com/pageza/recipe-book-api-v2/grpc/recipe"
	grpcUser "github.com/pageza/recipe-book-api-v2/grpc/user"

	"github.com/pageza/recipe-book-api-v2/internal/service"
)

// StartGRPCServer registers multiple gRPC services in a single gRPC server
func StartGRPCServer(userSvc service.UserService, recipeSvc service.RecipeService, notificationSvc service.NotificationService) error {
	lis, err := net.Listen("tcp", ":50051") // Change port as needed
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, grpcUser.NewServer(userSvc))
	pb.RegisterRecipeServiceServer(grpcServer, grpcRecipe.NewServer(recipeSvc))
	pb.RegisterNotificationServiceServer(grpcServer, grpcNotification.NewServer(&notificationSvc))

	// Enable gRPC Reflection (Required for grpcurl)
	reflection.Register(grpcServer)

	log.Println("gRPC server started on port 50051")
	return grpcServer.Serve(lis)
}
