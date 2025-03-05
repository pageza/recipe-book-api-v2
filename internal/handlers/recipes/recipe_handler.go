package recipes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/models"
)

// RecipeService defines the interface for recipe operations.
type RecipeService interface {
	CreateRecipe(recipe *models.Recipe) error
	GetRecipe(recipeID string) (*models.Recipe, error)
	ListRecipes() ([]*models.Recipe, error)
	ResolveRecipeQuery(query string) (*models.RecipeQueryResponse, error)
	// Additional service methods can be defined here.
}

// RecipeHandler handles HTTP requests related to recipes.
type RecipeHandler struct {
	service RecipeService
}

// NewRecipeHandler returns a new instance of RecipeHandler.
func NewRecipeHandler(svc RecipeService) *RecipeHandler {
	return &RecipeHandler{service: svc}
}

// @Summary Query Recipes
// @Description Queries recipes matching the search criteria, and creates one if needed.
// @Tags Recipes
// @Accept  json
// @Produce json
// @Param query body object{query=string} true "Recipe query request"
// @Success 200 {object} models.RecipeQueryResponse "Query result"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /query [post]
func (h *RecipeHandler) Query(c *gin.Context) {
	var payload struct {
		Query string `json:"query" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	resolverResp, err := h.service.ResolveRecipeQuery(payload.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resolverResp)
}

// Create handles the creation of a new recipe.
func (h *RecipeHandler) Create(c *gin.Context) {
	var payload models.Recipe
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recipe title is required"})
		return
	}
	if err := h.service.CreateRecipe(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, payload)
}

// Get retrieves a recipe by its ID.
func (h *RecipeHandler) Get(c *gin.Context) {
	recipeID := c.Param("id")
	recipe, err := h.service.GetRecipe(recipeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

// List returns all recipes.
func (h *RecipeHandler) List(c *gin.Context) {
	recipes, err := h.service.ListRecipes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list recipes"})
		return
	}
	c.JSON(http.StatusOK, recipes)
}
