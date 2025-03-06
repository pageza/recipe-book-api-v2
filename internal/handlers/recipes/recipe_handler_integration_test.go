package recipes_test

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	recipes "github.com/pageza/recipe-book-api-v2/internal/handlers/recipes"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
	"github.com/pageza/recipe-book-api-v2/internal/service"

	// generated proto package for recipes
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"gorm.io/gorm"

	pb "github.com/pageza/recipe-book-api-v2/proto/proto"
)

var testDB *gorm.DB // our test DB connection
var grpcClient pb.RecipeServiceClient
var router *gin.Engine

// TestMain sets up an in-memory SQLite database and spawns an in-process gRPC server.
func TestMain(m *testing.M) {
	// Use SQLite in-memory database for tests.
	os.Setenv("TEST_DB_DRIVER", "sqlite")
	log.Println("TEST_DB_DRIVER set to sqlite.")

	var err error
	// Connect to test database.
	testDB, err = repository.ConnectTestDB()
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}

	// Auto-migrate the Recipe model.
	if err = testDB.AutoMigrate(&models.Recipe{}); err != nil {
		log.Fatalf("failed to auto-migrate recipes table: %v", err)
	}
	log.Println("Auto-migration complete.")

	// Start an in-process gRPC server on a random free port.
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("failed to open test listener: %v", err)
	}
	log.Printf("gRPC listener opened on %s", lis.Addr())

	grpcServer := grpc.NewServer()

	// Instead of using the full service implementation, register the dummy server.
	pb.RegisterRecipeServiceServer(grpcServer, recipes.NewDummyRecipeServer())

	// Start the gRPC server in a separate goroutine.
	go func() {
		log.Println("Starting gRPC server (dummy)...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server terminated: %v", err)
		}
	}()
	log.Println("gRPC server goroutine started; waiting 500ms for startup...")
	time.Sleep(500 * time.Millisecond)

	// Dial the in-process gRPC server without blocking.
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial in-process gRPC server: %v", err)
	}

	// Poll until the connection becomes READY.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for {
		state := conn.GetState()
		log.Printf("Current connection state: %v", state)
		if state == connectivity.Ready {
			log.Println("Connection is READY.")
			break
		}
		if !conn.WaitForStateChange(ctx, state) {
			log.Fatalf("Timeout waiting for connection state to change from %v", state)
		}
	}
	grpcClient = pb.NewRecipeServiceClient(conn)
	log.Println("gRPC client setup complete.")

	// Initialize the router for HTTP integration tests using the local helper.
	recipeRepo := repository.NewRecipeRepository(testDB)
	recipeSvc := service.NewRecipeService(recipeRepo)
	router = setupRouter(recipeSvc)

	// Run tests.
	code := m.Run()

	grpcServer.GracefulStop()
	log.Println("gRPC server gracefully stopped.")
	os.Exit(code)
}

func TestGetRecipeGRPC(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetRecipeRequest{
		RecipeId: "test-recipe-1",
	}
	resp, err := grpcClient.GetRecipe(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// Additional assertions based on the response can be made here.
}

func TestQueryRecipeGRPC(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.RecipeQueryRequest{
		Query:  "vegan",
		UserId: "user-123",
		Filter: "dinner",
		Page:   1,
		Limit:  10,
	}
	resp, err := grpcClient.QueryRecipe(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// Validate the contents of resp as needed.
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
	// Clear the recipes table to avoid leftover data.
	if err := testDB.Exec("DELETE FROM recipes").Error; err != nil {
		t.Fatalf("failed to clear recipes table: %v", err)
	}

	// Seed a recipe with UserID "myUser".
	recipe := models.Recipe{
		ID:        "r123",
		UserID:    "myUser", // Correct field name
		Title:     "Test Recipe for myUser",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := testDB.Create(&recipe).Error; err != nil {
		t.Fatalf("failed to seed test recipe: %v", err)
	}

	// Issue GET request to the recipes endpoint with user_id=myUser.
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/recipes?user_id=myUser", nil)
	router.ServeHTTP(w, req)

	// Parse the JSON response.
	var resp struct {
		Recipes []struct {
			UserID string `json:"user_id"`
		} `json:"recipes"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(resp.Recipes) == 0 {
		t.Fatalf("expected at least one recipe")
	}

	// Assert that the returned recipe's user_id matches the queried user_id.
	assert.Equal(t, "myUser", resp.Recipes[0].UserID, "Recipe user_id should match queried user_id")
}

func TestQueryOtherUsersRecipes(t *testing.T) {
	// Clear the recipes table.
	if err := testDB.Exec("DELETE FROM recipes").Error; err != nil {
		t.Fatalf("failed to clear recipes table: %v", err)
	}

	// Seed a recipe with UserID "otherUser".
	recipe := models.Recipe{
		ID:        "r456",
		UserID:    "otherUser", // Correct field name
		Title:     "Test Recipe for otherUser",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := testDB.Create(&recipe).Error; err != nil {
		t.Fatalf("failed to seed test recipe: %v", err)
	}

	// Issue GET request to the recipes endpoint with user_id=otherUser.
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/recipes?user_id=otherUser", nil)
	router.ServeHTTP(w, req)

	// Parse the JSON response.
	var resp struct {
		Recipes []struct {
			UserID string `json:"user_id"`
		} `json:"recipes"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(resp.Recipes) == 0 {
		t.Fatalf("expected at least one recipe")
	}

	// Assert that the returned recipe's user_id matches the queried user_id.
	assert.Equal(t, "otherUser", resp.Recipes[0].UserID, "Recipe user_id should match queried user_id")
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
