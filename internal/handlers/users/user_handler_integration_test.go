package users_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pageza/recipe-book-api-v2/proto/proto" // Generated gRPC client code
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// setupTestClient connects to the gRPC server using the address specified in the environment variable GRPC_DIAL_ADDRESS.
// If not set, it defaults to "grpc-server:50051".
func setupTestClient() proto.UserServiceClient {
	// Allow time for the server to be ready.
	time.Sleep(2 * time.Second)

	addr := os.Getenv("GRPC_DIAL_ADDRESS")
	if addr == "" {
		addr = "grpc-server:50051"
	}

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
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

	// Wait for the record to be committed.
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

func TestIntegration_DuplicateRegistration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	client := setupTestClient()

	// Generate unique email and username for this test.
	uniqueEmail := "inttestdup_" + uuid.New().String() + "@example.com"
	uniqueUsername := "inttestdup_" + uuid.New().String()

	// 1. First registration should succeed.
	regResp1, err := client.Register(context.Background(), &proto.CreateUserRequest{
		Email:       uniqueEmail,
		Username:    uniqueUsername,
		Password:    "duplicatepassword",
		Preferences: "{\"diet\":\"vegan\"}",
	})
	assert.NoError(t, err, "Expected no error on first registration")
	assert.NotEmpty(t, regResp1.UserId, "Expected userId in first registration response")

	// Wait for commit.
	time.Sleep(2 * time.Second)

	// 2. Second registration with the same email should fail.
	_, err = client.Register(context.Background(), &proto.CreateUserRequest{
		Email:       uniqueEmail,
		Username:    uniqueUsername,
		Password:    "duplicatepassword",
		Preferences: "{\"diet\":\"vegan\"}",
	})
	assert.Error(t, err, "Expected error on duplicate registration")
	// Check that the error message contains "user already exists".
	st, ok := status.FromError(err)
	assert.True(t, ok, "Expected gRPC status error")
	assert.Contains(t, st.Message(), "user already exists", "Expected duplicate registration error message")
}

func TestIntegration_RegisterEmptyEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	client := setupTestClient()

	// Create a registration request with an empty email.
	regResp, err := client.Register(context.Background(), &proto.CreateUserRequest{
		Email:       "",
		Username:    "malformedUser",
		Password:    "somepassword",
		Preferences: "{\"diet\":\"vegan\"}",
	})
	// In our environment, an empty email might already be registered, so we expect an error.
	assert.Error(t, err, "Expected error for registration with empty email")
	st, ok := status.FromError(err)
	assert.True(t, ok, "Expected gRPC status error")
	// Adjust the expected substring as per your service's response.
	assert.Contains(t, st.Message(), "email cannot be empty", "Expected error message indicating missing email")
	// Optionally, if a response is returned, ensure that the Email field in the response is empty.
	if regResp != nil {
		// Since CreateUserResponse doesn't have an Email field, we check for UserId only.
		assert.NotEmpty(t, regResp.UserId, "Expected userId even for empty email registration")
	}
}
