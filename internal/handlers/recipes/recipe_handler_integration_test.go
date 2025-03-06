package recipes_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	recipes "github.com/pageza/recipe-book-api-v2/internal/handlers/recipes"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/repository" // make sure repository package is imported
	// generated proto package for recipes
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	pb "github.com/pageza/recipe-book-api-v2/proto/proto"
)

var testDB *gorm.DB // our test DB connection
var grpcClient pb.RecipeServiceClient

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
	grpcClient = pb.NewRecipeServiceClient(conn)

	// Run tests.
	code := m.Run()

	// Optionally, clean up the test database here.
	// repository.CleanupTestDB(testDB)

	os.Exit(code)
}

func TestGetRecipeGRPC(t *testing.T) {
	// Replace with the correct address of your test gRPC server
	address := "localhost:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewRecipeServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetRecipeRequest{
		RecipeId: "test-recipe-1",
	}
	resp, err := client.GetRecipe(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// Additional assertions based on the expected data from the GetRecipe RPC.
}

func TestQueryRecipeGRPC(t *testing.T) {
	address := "localhost:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewRecipeServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.RecipeQueryRequest{
		Query:  "vegan",
		UserId: "user-123",
		Filter: "dinner",
		Page:   1,
		Limit:  10,
	}
	resp, err := client.QueryRecipe(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// Validate the response, e.g.
	// assert.Equal(t, int32(1), resp.Page)
	// assert.Equal(t, int32(10), resp.Limit)
	// assert.NotNil(t, resp.Recipes)
}

// mockRecipeService implements recipes.RecipeService for testing.
type mockRecipeService struct{}

func (m *mockRecipeService) GetRecipe(recipeID string) (*models.Recipe, error) {
	return &models.Recipe{
		ID:                recipeID,
		Title:             "Test Recipe",
		Ingredients:       "Ingredient1, Ingredient2",
		Steps:             "Step1, Step2",
		NutritionalInfo:   "Calories: 100",
		AllergyDisclaimer: "None",
		Appliances:        "Oven, Stove",
		CreatedAt:         time.Unix(1630000000, 0),
		UpdatedAt:         time.Unix(1630000000, 0),
		UserID:            "user123",
	}, nil
}

func (m *mockRecipeService) QueryRecipes(req *models.RecipeQueryRequest) (*models.RecipeQueryResponse, error) {
	// For testing, return two dummy recipes.
	recipesList := []*models.Recipe{
		{
			ID:                "r1",
			Title:             "Recipe One",
			Ingredients:       "Ing1, Ing2",
			Steps:             "Step1, Step2",
			NutritionalInfo:   "Info1",
			AllergyDisclaimer: "None",
			Appliances:        "Microwave, Oven",
			CreatedAt:         time.Unix(1630000000, 0),
			UpdatedAt:         time.Unix(1630000000, 0),
			UserID:            "user123",
		},
		{
			ID:                "r2",
			Title:             "Recipe Two",
			Ingredients:       "IngA, IngB",
			Steps:             "StepA, StepB",
			NutritionalInfo:   "Info2",
			AllergyDisclaimer: "None",
			Appliances:        "Stove",
			CreatedAt:         time.Unix(1630000000, 0),
			UpdatedAt:         time.Unix(1630000000, 0),
			UserID:            "user123",
		},
	}

	return &models.RecipeQueryResponse{
		Recipes: recipesList,
		Page:    req.Page,
		Limit:   req.Limit,
		Total:   2,
	}, nil
}

// setupRouter initializes a Gin router with the RecipeHandler routes.
func setupRouter(service recipes.RecipeService) *gin.Engine {
	router := gin.Default()
	handler := recipes.NewRecipeHandler(service)
	router.GET("/recipes", handler.Query)
	router.GET("/recipes/:id", handler.Get)
	return router
}

func TestGetRecipe(t *testing.T) {
	// Test scenario: Retrieve a recipe by its ID.
	router := setupRouter(&mockRecipeService{})
	req, _ := http.NewRequest("GET", "/recipes/r123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "Expected HTTP status 200 for GetRecipe")

	var recipe models.Recipe
	err := json.Unmarshal(w.Body.Bytes(), &recipe)
	assert.NoError(t, err, "Expected no error unmarshalling response")
	assert.Equal(t, "r123", recipe.ID, "Recipe ID should match")
	assert.Equal(t, "Test Recipe", recipe.Title, "Recipe title should be 'Test Recipe'")
}

func TestQueryMyRecipes(t *testing.T) {
	// Test scenario: User wants to view their own recipes.
	router := setupRouter(&mockRecipeService{})
	req, _ := http.NewRequest("GET", "/recipes?user_id=myUser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "Expected HTTP 200 for Query My Recipes")

	var response models.RecipeQueryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected no error unmarshalling query response")
	assert.Equal(t, 2, len(response.Recipes), "Should return two recipes for the user")
	for _, recipe := range response.Recipes {
		assert.Equal(t, "myUser", recipe.UserID, "Recipe user_id should match queried user_id")
	}
}

func TestQueryOtherUsersRecipes(t *testing.T) {
	// Test scenario: Retrieve recipes created by another user.
	router := setupRouter(&mockRecipeService{})
	req, _ := http.NewRequest("GET", "/recipes?user_id=otherUser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "Expected HTTP 200 for Query Other User's Recipes")

	var response models.RecipeQueryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected no error unmarshalling response")
	assert.Equal(t, 2, len(response.Recipes), "Should return two recipes for another user")
	for _, recipe := range response.Recipes {
		assert.Equal(t, "otherUser", recipe.UserID, "Recipe user_id should match queried user_id")
	}
}

func TestQueryByCuisineOrDiet(t *testing.T) {
	router := setupRouter(&mockRecipeService{})
	req, _ := http.NewRequest("GET", "/recipes?query=vegan&page=1&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "Expected HTTP 200 for Query by Cuisine/Diet")

	var response models.RecipeQueryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected no error unmarshalling response")
	assert.Equal(t, 2, len(response.Recipes), "Should return two recipes for advanced query")
	assert.Equal(t, 1, response.Page, "Page should be 1")
	assert.Equal(t, 10, response.Limit, "Limit should be 10")
}

func TestQueryCombinedFilters(t *testing.T) {
	router := setupRouter(&mockRecipeService{})
	queryParams := "query=quick&user_id=testUser&filter=Indian&page=2&limit=5"
	req, _ := http.NewRequest("GET", "/recipes?"+queryParams, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "Expected HTTP 200 for combined filters query")

	var response models.RecipeQueryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected no error unmarshalling response")
	// In our mock, the response always returns two recipes.
	assert.Equal(t, 2, len(response.Recipes), "Should return two recipes for combined query")
	page, _ := strconv.Atoi("2")
	limit, _ := strconv.Atoi("5")
	assert.Equal(t, page, response.Page, "Page should be echoed as 2")
	assert.Equal(t, limit, response.Limit, "Limit should be echoed as 5")
}
