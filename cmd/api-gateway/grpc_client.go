package main

import (
	"context"

	pb "github.com/pageza/recipe-book-api-v2/proto/proto"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	UserServiceClient pb.UserServiceClient
}

// NewGRPCClient initializes gRPC clients
func NewGRPCClient() (*GRPCClient, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // Use TLS in production!
	if err != nil {
		return nil, err
	}

	return &GRPCClient{
		UserServiceClient: pb.NewUserServiceClient(conn),
	}, nil
}

// Example: Call UserService Register
func (c *GRPCClient) RegisterUser(username, email, password string) error {
	req := &pb.CreateUserRequest{
		Username: username,
		Email:    email,
		Password: password,
	}
	_, err := c.UserServiceClient.Register(context.Background(), req)
	return err
}
