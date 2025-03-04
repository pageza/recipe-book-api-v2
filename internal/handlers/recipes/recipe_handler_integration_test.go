package recipes_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pageza/recipe-book-api-v2/internal/handlers/recipes"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/repository" // make sure repository package is imported
	"github.com/pageza/recipe-book-api-v2/proto/proto"         // generated proto package for recipes
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var testDB *gorm.DB // our test DB connection
var grpcClient proto.RecipeServiceClient

// TestMain sets up the test database and migrates the Recipe model before running tests.
func TestMain(m *testing.M) {
	var err error
	// Connect to your test database.
	testDB, err = repository.ConnectTestDB() // ensure you have this helper to connect to your test DB
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}

	// Auto-migrate the Recipe model.
	err = testDB.AutoMigrate(&models.Recipe{})
	if err != nil {
		log.Fatalf("failed to auto-migrate recipes table: %v", err)
	}

	// Optionally, wait a bit to ensure the migration is done.
	time.Sleep(1 * time.Second)

	// Setup gRPC client connection for integration tests.
	// For example, using an environment variable or a default address.
	addr := os.Getenv("GRPC_DIAL_ADDRESS")
	if addr == "" {
		addr = "grpc-server:50051" // adjust if needed
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	grpcClient = proto.NewRecipeServiceClient(conn)

	// Run tests.
	code := m.Run()

	// Optionally, clean up the test database here.
	// repository.CleanupTestDB(testDB)

	os.Exit(code)
}

func TestIntegration_CreateAndGetRecipe(t *testing.T) {
	// Read the gRPC server address from the environment, just like in user_handler_integration_test.go.
	grpcServerAddr := os.Getenv("GRPC_SERVER_ADDR")
	if grpcServerAddr == "" {
		grpcServerAddr = "grpc-server:50051"
	}

	conn, err := grpc.Dial(grpcServerAddr, grpc.WithInsecure())
	assert.NoError(t, err, "Expected to connect to gRPC server")
	defer conn.Close()

	client := proto.NewRecipeServiceClient(conn)

	// Example: create a new recipe.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.CreateRecipe(ctx, &proto.CreateRecipeRequest{
		Title:       "Healthy Chicken Salad",
		Ingredients: `{"items": ["chicken", "lettuce", "tomato", "cucumber"]}`,
		Steps:       `{"steps": ["Grill chicken", "Chop veggies", "Mix together"]}`,
		// Note: our new Recipe model has additional fields (NutritionalInfo, AllergyDisclaimer, Appliances)
		// but the current proto definition doesn't include these.
		// We assume that later you'll update the proto, so for now these are omitted.
	})
	assert.NoError(t, err, "Expected no error during recipe creation")
	assert.NotEmpty(t, resp.RecipeId, "Expected a non-empty recipeId")

	// Wait briefly for the recipe to be available.
	time.Sleep(1 * time.Second)

	// Retrieve the recipe via gRPC.
	getReq := &proto.GetRecipeRequest{
		RecipeId: resp.RecipeId,
	}
	getResp, err := client.GetRecipe(ctx, getReq)
	assert.NoError(t, err, "Expected no error during recipe retrieval")
	assert.Equal(t, "Healthy Chicken Salad", getResp.Title, "Recipe title should match")
	assert.Equal(t, `{"items": ["chicken", "lettuce", "tomato", "cucumber"]}`, getResp.Ingredients, "Ingredients should match")
	assert.Equal(t, `{"steps": ["Grill chicken", "Chop veggies", "Mix together"]}`, getResp.Steps, "Steps should match")
}

func TestIntegration_ListRecipes(t *testing.T) {
	// This test assumes that there is at least one recipe in the database from previous tests.
	listReq := &proto.ListRecipesRequest{}
	listResp, err := grpcClient.ListRecipes(context.Background(), listReq)
	assert.NoError(t, err, "Expected no error during listing recipes")
	assert.Greater(t, len(listResp.Recipes), 0, "Expected at least one recipe in the list")
}

func TestIntegration_QueryRecipe(t *testing.T) {
	// Since the query functionality is part of our application logic (using RAG and so on),
	// we'll assume that the service layer wraps the proto.GetRecipe functionality for now.
	// For this test, we simulate a query by first creating a recipe and then querying it.
	uniqueTitle := "Vegan Delight " + uuid.New().String()
	createReq := &proto.CreateRecipeRequest{
		Title:       uniqueTitle,
		Ingredients: `{"items": ["tofu", "spinach", "quinoa"]}`,
		Steps:       `{"steps": ["Cook quinoa", "Sauté tofu", "Mix with spinach"]}`,
	}
	createResp, err := grpcClient.CreateRecipe(context.Background(), createReq)
	assert.NoError(t, err, "Expected no error during recipe creation for query")

	// Here, we assume that the query returns a primary recipe and alternatives.
	// For now, we simulate by using GetRecipe.
	getReq := &proto.GetRecipeRequest{
		RecipeId: createResp.RecipeId,
	}
	getResp, err := grpcClient.GetRecipe(context.Background(), getReq)
	assert.NoError(t, err, "Expected no error during recipe query simulation")
	// We simulate the RecipeQueryResponse here by asserting that we got the expected title.
	assert.Equal(t, uniqueTitle, getResp.Title, "Queried recipe title should match the created recipe")
}

// dummyRecipeService implements the minimal RecipeService interface needed for testing.
type dummyRecipeService struct{}

// ResolveRecipeQuery returns a fixed dummy recipe response.
func (d *dummyRecipeService) ResolveRecipeQuery(query string) (*models.RecipeQueryResponse, error) {
	return &models.RecipeQueryResponse{
		Recipes: []*models.Recipe{
			{
				ID:                "dummy-id",
				Title:             query + " - Dummy Generated Recipe",
				Ingredients:       `["dummy ingredient"]`,
				Steps:             `["dummy step"]`,
				NutritionalInfo:   "{}",
				AllergyDisclaimer: "none",
				Appliances:        "[]",
			},
		},
	}, nil
}

func TestQueryRecipe(t *testing.T) {
	// Set Gin to test mode.
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Create a dummy service and attach it to the handler.
	dummySvc := &dummyRecipeService{}
	handler := recipes.NewRecipeHandler(dummySvc)

	// Register the /query endpoint.
	router.POST("/query", handler.Query)

	// Prepare a test request.
	reqBody := `{"query": "Test Recipe"}`
	req, err := http.NewRequest("POST", "/query", bytes.NewBufferString(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request.
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response status code.
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response.
	var resp models.RecipeQueryResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Verify the response.
	assert.Len(t, resp.Recipes, 1, "Expected a single recipe in the response")
	assert.Equal(t, "dummy-id", resp.Recipes[0].ID, "Recipe ID should match the dummy one")
	assert.Contains(t, resp.Recipes[0].Title, "Test Recipe", "Title should include the query string")
}
