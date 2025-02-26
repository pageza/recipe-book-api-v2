package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"
	pb "github.com/pageza/recipe-book-api-v2/proto/proto"
	"github.com/pkg/errors"
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
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: "",
		Preferences:  req.Preferences,
	}
	user.ID = uuid.New().String()

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash password")
	}
	user.PasswordHash = hashed

	if err := s.svc.Register(user); err != nil {
		return nil, errors.Wrap(err, "failed to register user")
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
		return nil, errors.Wrap(err, "login failed")
	}

	token, err := utils.GenerateJWT(user.ID, "user", []string{"read:profile"}, "testsecret")
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate token")
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
		return nil, errors.Wrap(err, "user not found")
	}

	return &pb.GetProfileResponse{
		Username:    user.Username,
		Email:       user.Email,
		Preferences: user.Preferences,
	}, nil
}
