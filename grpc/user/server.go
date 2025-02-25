package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"
	pb "github.com/pageza/recipe-book-api-v2/proto/proto"
)

// Server implements the gRPC UserService.
type Server struct {
	pb.UnimplementedUserServiceServer
	svc service.UserService
}

// NewServer creates a new User gRPC server.
func NewServer(svc service.UserService) *Server {
	return &Server{svc: svc}
}

// Register implements the Register RPC.
func (s *Server) Register(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Convert proto request to internal model
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: "", // Will be generated below
		Preferences:  req.Preferences,
	}

	// Generate a user ID
	user.ID = uuid.New().String()

	// Hash the password
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}
	user.PasswordHash = hashed

	// Call business logic
	if err := s.svc.Register(user); err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	return &pb.CreateUserResponse{
		UserId:  user.ID,
		Message: "User registered successfully",
	}, nil
}

// Login implements the Login RPC.
func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.svc.Login(req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("login failed: %v", err)
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, "your-secret") // Replace with actual config-based secret
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &pb.LoginResponse{
		UserId: user.ID,
		Token:  token,
	}, nil
}

// GetProfile implements the GetProfile RPC.
func (s *Server) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	user, err := s.svc.GetProfile(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	return &pb.GetProfileResponse{
		Username:    user.Username,
		Email:       user.Email,
		Preferences: user.Preferences,
	}, nil
}
