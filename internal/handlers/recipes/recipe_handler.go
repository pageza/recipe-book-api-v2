package recipes

import (
	"net/http"
	"strconv"
	"strings"

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

// Query handles GET requests for querying recipes.
// This unified endpoint interprets the parameters to support:
//   - Listing recipes by user (if query is empty but user_id is provided).
//   - Advanced search queries (if query is non-empty).
//
// Endpoint: GET /recipes
func (h *RecipeHandler) Query(c *gin.Context) {
	queryParam := strings.TrimSpace(c.DefaultQuery("query", ""))
	userID := c.Query("user_id")
	filter := c.Query("filter")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// Convert pagination parameters to integers.
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Construct the RecipeQueryRequest from the URL parameters.
	req := &models.RecipeQueryRequest{
		Query:  queryParam,
		UserID: userID,
		Filter: filter,
		Page:   page,
		Limit:  limit,
	}

	// Delegate the query to the RecipeService.
	resp, err := h.service.QueryRecipes(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process query: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
