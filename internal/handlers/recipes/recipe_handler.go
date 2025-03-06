package recipes

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/models"
)

// RecipeService defines the interface for recipe operations.
type RecipeService interface {
	// GetRecipe retrieves a recipe by its ID.
	GetRecipe(recipeID string) (*models.Recipe, error)
	// QueryRecipes processes query requests for recipes.
	QueryRecipes(req *models.RecipeQueryRequest) (*models.RecipeQueryResponse, error)
}

// RecipeHandler handles HTTP requests related to recipes.
type RecipeHandler struct {
	service RecipeService
}

// NewRecipeHandler constructs a new RecipeHandler with the given RecipeService.
func NewRecipeHandler(service RecipeService) *RecipeHandler {
	return &RecipeHandler{service: service}
}

// Get handles GET requests to retrieve a single recipe by its ID.
// Endpoint: GET /recipes/:id
func (h *RecipeHandler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recipe id is required"})
		return
	}
	recipe, err := h.service.GetRecipe(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

// Query handles POST requests to /recipe/query.
// It now binds JSON from the request body (instead of reading URL query parameters)
// and forwards the {"query": "..."} payload to the resolver microservice.
func (h *RecipeHandler) Query(c *gin.Context) {
	log.Println("Main App: Received POST /recipe/query")

	// Bind the incoming JSON payload into a RecipeQueryRequest.
	var req models.RecipeQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Main App: Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	log.Printf("Main App: Bound JSON request: %+v", req)

	// Forward the query to the resolver microservice.
	// Use the Docker service name (e.g., "resolver") based on the docker-compose configuration.
	resolverURL := "http://resolver:3000/resolve"
	// Prepare a payload that the resolver expects: {"query": "..." }
	payload, err := json.Marshal(map[string]string{
		"query": req.Query,
	})
	if err != nil {
		log.Printf("Main App: Error marshalling payload: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	log.Printf("Main App: Sending payload to resolver: %s", payload)

	// Call the resolver microservice.
	resp, err := http.Post(resolverURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Printf("Main App: Error calling resolver service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	defer resp.Body.Close()

	// Read the resolver's response.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Main App: Error reading resolver response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	log.Printf("Main App: Received response from resolver: %s", body)

	// Define a local type that matches the structure the resolver returns.
	type ResolverResponse struct {
		PrimaryRecipe      models.Recipe   `json:"primary_recipe"`
		AlternativeRecipes []models.Recipe `json:"alternative_recipes"`
	}
	var resolverResp ResolverResponse
	if err := json.Unmarshal(body, &resolverResp); err != nil {
		log.Printf("Main App: Error unmarshalling resolver response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	log.Printf("Main App: Resolver returned primary: %+v, alternatives: %+v", resolverResp.PrimaryRecipe, resolverResp.AlternativeRecipes)

	// Return the resolver's response to the client.
	c.JSON(http.StatusOK, resolverResp)
}
