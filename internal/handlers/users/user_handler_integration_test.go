package users_test

import (
	"context"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pageza/recipe-book-api-v2/proto/proto" // Generated gRPC client code
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func setupTestClient() proto.UserServiceClient {
	// Allow some time for the server to be ready.
	time.Sleep(2 * time.Second)

	// Connect using the Docker service name.
	conn, err := grpc.Dial("grpc-server:50051", grpc.WithInsecure())
	if err != nil {
		panic("Failed to connect to gRPC server: " + err.Error())
	}
	return proto.NewUserServiceClient(conn)
}

func TestIntegration_RegisterAndLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	client := setupTestClient()

	// Generate unique email and username to avoid duplicates.
	uniqueEmail := "inttestuser_" + uuid.New().String() + "@example.com"
	uniqueUsername := "inttestuser_" + uuid.New().String()

	// 1. Register via gRPC.
	regResp, err := client.Register(context.Background(), &proto.CreateUserRequest{
		Email:       uniqueEmail,
		Username:    uniqueUsername,
		Password:    "inttestpassword", // Plain password as defined in proto.
		Preferences: "{\"diet\":\"vegan\"}",
	})
	assert.NoError(t, err, "Expected no error during registration")
	assert.NotEmpty(t, regResp.UserId, "Expected userId in registration response")

	// Wait for the record to be fully committed.
	time.Sleep(2 * time.Second)

	// 2. Login via gRPC.
	loginResp, err := client.Login(context.Background(), &proto.LoginRequest{
		Email:    uniqueEmail,
		Password: "inttestpassword",
	})
	assert.NoError(t, err, "Expected no error during login")
	assert.NotEmpty(t, loginResp.Token, "Expected token in login response")
	assert.NotEmpty(t, loginResp.UserId, "Expected userId in login response")

	// 3. Get Profile via gRPC.
	profileResp, err := client.GetProfile(context.Background(), &proto.GetProfileRequest{
		UserId: loginResp.UserId,
	})
	assert.NoError(t, err, "Expected no error during profile fetch")
	assert.Equal(t, uniqueEmail, profileResp.Email, "Profile email should match")
	assert.Equal(t, uniqueUsername, profileResp.Username, "Profile username should match")
}

func TestIntegration_InvalidLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	client := setupTestClient()

	// Generate unique email and username for this test.
	uniqueEmail := "inttestuser2_" + uuid.New().String() + "@example.com"
	uniqueUsername := "inttestuser2_" + uuid.New().String()

	// 1. Register a user for testing invalid login.
	_, err := client.Register(context.Background(), &proto.CreateUserRequest{
		Email:       uniqueEmail,
		Username:    uniqueUsername,
		Password:    "validpassword",
		Preferences: "{\"diet\":\"vegetarian\"}",
	})
	assert.NoError(t, err)

	// Wait briefly for the record to be available.
	time.Sleep(2 * time.Second)

	// 2. Attempt to log in with the wrong password.
	loginResp, err := client.Login(context.Background(), &proto.LoginRequest{
		Email:    uniqueEmail,
		Password: "wrongpassword",
	})
	assert.Error(t, err, "Expected error during login with incorrect password")
	// Since an error is expected, loginResp should be nil.
	assert.Nil(t, loginResp, "Expected login response to be nil when login fails")
}
