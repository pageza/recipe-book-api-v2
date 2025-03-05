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
	"github.com/pageza/recipe-book-api-v2/internal/repository" // ensure repository package is imported
	"github.com/pageza/recipe-book-api-v2/internal/service"
	pb "github.com/pageza/recipe-book-api-v2/proto/proto" // generated proto package for recipes
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"gorm.io/gorm"
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

func TestIntegration_CreateAndGetRecipe(t *testing.T) {
	// Read the gRPC server address from the environment.
	grpcServerAddr := os.Getenv("GRPC_SERVER_ADDR")
	if grpcServerAddr == "" {
		grpcServerAddr = "grpc-server:50051"
	}

	conn, err := grpc.Dial(grpcServerAddr, grpc.WithInsecure())
	assert.NoError(t, err, "Expected to connect to gRPC server")
	defer conn.Close()

	client := pb.NewRecipeServiceClient(conn)

	// Create a new recipe via gRPC.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createResp, err := client.CreateRecipe(ctx, &pb.CreateRecipeRequest{
		Title:       "Healthy Chicken Salad",
		Ingredients: `{"items": ["chicken", "lettuce", "tomato", "cucumber"]}`,
		Steps:       `{"steps": ["Grill chicken", "Chop veggies", "Mix together"]}`,
	})
	assert.NoError(t, err, "Expected no error during recipe creation")
	assert.NotEmpty(t, createResp.RecipeId, "Expected a non-empty recipeId")

	// Wait briefly for the recipe to be available.
	time.Sleep(1 * time.Second)

	// Retrieve the recipe via gRPC.
	getResp, err := client.GetRecipe(ctx, &pb.GetRecipeRequest{
		RecipeId: createResp.RecipeId,
	})
	assert.NoError(t, err, "Expected no error during recipe retrieval")
	assert.Equal(t, "Healthy Chicken Salad", getResp.Title, "Recipe title should match")
	assert.Equal(t, `{"items": ["chicken", "lettuce", "tomato", "cucumber"]}`, getResp.Ingredients, "Ingredients should match")
	assert.Equal(t, `{"steps": ["Grill chicken", "Chop veggies", "Mix together"]}`, getResp.Steps, "Steps should match")
}

func TestIntegration_ListRecipes(t *testing.T) {
	listResp, err := grpcClient.ListRecipes(context.Background(), &pb.ListRecipesRequest{})
	assert.NoError(t, err, "Expected no error during listing recipes")
	assert.Greater(t, len(listResp.Recipes), 0, "Expected at least one recipe in the list")
}

func TestIntegration_QueryRecipe(t *testing.T) {
	// Simulate a query by first creating a unique recipe.
	uniqueTitle := "Vegan Delight " + uuid.New().String()
	createReq := &pb.CreateRecipeRequest{
		Title:       uniqueTitle,
		Ingredients: `{"items": ["tofu", "spinach", "quinoa"]}`,
		Steps:       `{"steps": ["Cook quinoa", "Sauté tofu", "Mix with spinach"]}`,
	}
	createResp, err := grpcClient.CreateRecipe(context.Background(), createReq)
	assert.NoError(t, err, "Expected no error during recipe creation for query")

	// Here we simulate the query functionality by retrieving the created recipe.
	getResp, err := grpcClient.GetRecipe(context.Background(), &pb.GetRecipeRequest{
		RecipeId: createResp.RecipeId,
	})
	assert.NoError(t, err, "Expected no error during recipe query simulation")
	assert.Equal(t, uniqueTitle, getResp.Title, "Queried recipe title should match the created recipe")
}

// --- HTTP Integration Tests for Recipes ---

// dummyError is used to simulate errors.
type dummyError struct {
	msg string
}

func (e *dummyError) Error() string {
	return e.msg
}

// dummyRecipeService implements RecipeService for GET and LIST endpoint testing.
type dummyRecipeService struct {
	recipes map[string]*models.Recipe
}

func (d *dummyRecipeService) CreateRecipe(recipe *models.Recipe) error {
	if d.recipes == nil {
		d.recipes = make(map[string]*models.Recipe)
	}
	if recipe.ID == "" {
		recipe.ID = "dummy-" + recipe.Title
	}
	d.recipes[recipe.ID] = recipe
	return nil
}

func (d *dummyRecipeService) GetRecipe(recipeID string) (*models.Recipe, error) {
	if recipe, ok := d.recipes[recipeID]; ok {
		return recipe, nil
	}
	return nil, &dummyError{msg: "recipe not found"}
}

func (d *dummyRecipeService) ListRecipes() ([]*models.Recipe, error) {
	var list []*models.Recipe
	for _, r := range d.recipes {
		list = append(list, r)
	}
	return list, nil
}

func (d *dummyRecipeService) ResolveRecipeQuery(query string) (*models.RecipeQueryResponse, error) {
	// Not needed for these tests.
	return nil, nil
}

// setupRouter initializes a Gin engine using dummyRecipeService for GET and LIST endpoints.
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	dummySvc := &dummyRecipeService{
		recipes: map[string]*models.Recipe{
			"dummy-Pasta": {
				ID:                "dummy-Pasta",
				Title:             "Pasta",
				Ingredients:       `["pasta", "tomato sauce"]`,
				Steps:             `["boil pasta", "add sauce"]`,
				NutritionalInfo:   "{}",
				AllergyDisclaimer: "",
				Appliances:        "[]",
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
		},
	}
	handler := recipes.NewRecipeHandler(dummySvc)
	router := gin.Default()
	router.GET("/recipes/:id", handler.Get)
	router.GET("/recipes", handler.List)
	return router
}

func TestGetRecipe(t *testing.T) {
	router := setupRouter()
	req, err := http.NewRequest("GET", "/recipes/dummy-Pasta", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var recipe models.Recipe
	err = json.Unmarshal(w.Body.Bytes(), &recipe)
	assert.NoError(t, err)
	assert.Equal(t, "dummy-Pasta", recipe.ID)
	assert.Equal(t, "Pasta", recipe.Title)
}

func TestListRecipes(t *testing.T) {
	router := setupRouter()
	req, err := http.NewRequest("GET", "/recipes", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var recipes []models.Recipe
	err = json.Unmarshal(w.Body.Bytes(), &recipes)
	assert.NoError(t, err)
	// Expect exactly one recipe seeded in the dummy service.
	assert.Len(t, recipes, 1)
}

// dummyRecipeServiceSimple implements a minimal RecipeService for POST Query and Create endpoint testing.
type dummyRecipeServiceSimple struct{}

func (d *dummyRecipeServiceSimple) ResolveRecipeQuery(query string) (*models.RecipeQueryResponse, error) {
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
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
		},
	}, nil
}

func (d *dummyRecipeServiceSimple) CreateRecipe(recipe *models.Recipe) error {
	recipe.ID = "dummy-created-id"
	return nil
}

func (d *dummyRecipeServiceSimple) GetRecipe(recipeID string) (*models.Recipe, error) {
	return nil, nil
}

func (d *dummyRecipeServiceSimple) ListRecipes() ([]*models.Recipe, error) {
	return nil, nil
}

func TestQueryRecipe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	dummySvc := &dummyRecipeServiceSimple{}
	handler := recipes.NewRecipeHandler(dummySvc)

	router.POST("/query", handler.Query)

	reqBody := `{"query": "Test Recipe"}`
	req, err := http.NewRequest("POST", "/query", bytes.NewBufferString(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.RecipeQueryResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Len(t, resp.Recipes, 1)
	assert.Equal(t, "dummy-id", resp.Recipes[0].ID)
	assert.Contains(t, resp.Recipes[0].Title, "Test Recipe")
}

func TestCreateRecipe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	dummySvc := &dummyRecipeServiceSimple{}
	handler := recipes.NewRecipeHandler(dummySvc)

	router.POST("/create", handler.Create)

	testRecipe := models.Recipe{
		Title:             "New Recipe",
		Ingredients:       `["ingredient1", "ingredient2"]`,
		Steps:             `["step1", "step2"]`,
		NutritionalInfo:   "{}",
		AllergyDisclaimer: "none",
		Appliances:        "[]",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	reqBytes, err := json.Marshal(testRecipe)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(reqBytes))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdRecipe models.Recipe
	err = json.Unmarshal(w.Body.Bytes(), &createdRecipe)
	assert.NoError(t, err)
	assert.Equal(t, "dummy-created-id", createdRecipe.ID)
	assert.Equal(t, "New Recipe", createdRecipe.Title)
}

// fakeResolverResponse mirrors the expected resolver response contract.
type fakeResolverResponse struct {
	PrimaryRecipe      *models.Recipe   `json:"primary_recipe"`
	AlternativeRecipes []*models.Recipe `json:"alternative_recipes"`
}

func TestRecipeQueryIntegrationSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a fake resolver service endpoint that returns the legacy structure.
	resolverHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			PrimaryRecipe      *models.Recipe   `json:"primary_recipe"`
			AlternativeRecipes []*models.Recipe `json:"alternative_recipes"`
		}{
			PrimaryRecipe: &models.Recipe{
				ID:                "integration-test-id",
				Title:             "Integration Test Recipe",
				Ingredients:       `["ingredient1", "ingredient2"]`,
				Steps:             `["step1", "step2"]`,
				NutritionalInfo:   "{}",
				AllergyDisclaimer: "none",
				Appliances:        "[]",
			},
			AlternativeRecipes: []*models.Recipe{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	testResolverServer := httptest.NewServer(resolverHandler)
	defer testResolverServer.Close()

	// Override the resolver service URL to use our fake server.
	service.ResolverServiceURL = testResolverServer.URL

	router := gin.Default()
	// Inject the RecipeService that delegates the resolution to the external resolver.
	handler := recipes.NewRecipeHandler(service.NewRecipeService(nil))
	router.POST("/query", handler.Query)

	reqBody := `{"query": "Integration Test Recipe"}`
	req, err := http.NewRequest("POST", "/query", bytes.NewBufferString(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	// The new contract returns the recipes in a "recipes" field.
	var queryResp struct {
		Recipes []*models.Recipe `json:"recipes"`
	}
	err = json.Unmarshal(recorder.Body.Bytes(), &queryResp)
	assert.NoError(t, err)
	assert.Len(t, queryResp.Recipes, 1)
	assert.Equal(t, "integration-test-id", queryResp.Recipes[0].ID)
	assert.Equal(t, "Integration Test Recipe", queryResp.Recipes[0].Title)
}
