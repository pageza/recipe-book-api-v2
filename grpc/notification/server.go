package notification

import (
	"context"
	"fmt"
	"log"

	"github.com/pageza/recipe-book-api-v2/internal/service"
	pb "github.com/pageza/recipe-book-api-v2/proto/proto"
)

// Server implements the gRPC NotificationService.
type Server struct {
	pb.UnimplementedNotificationServiceServer
	svc *service.NotificationService
}

// NewServer creates a new Notification gRPC server.
func NewServer(svc *service.NotificationService) *Server {
	return &Server{svc: svc}
}

// SendNotification implements the SendNotification RPC.
func (s *Server) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	err := s.svc.SendNotification(req.UserId, req.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to send notification: %v", err)
	}

	log.Println("Notification sent successfully")

	return &pb.SendNotificationResponse{
		Status: "Notification sent successfully",
	}, nil
}
