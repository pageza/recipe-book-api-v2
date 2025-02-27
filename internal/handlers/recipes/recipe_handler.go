package recipes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/models"
)

// RecipeService defines the interface for querying recipes.
type RecipeService interface {
	// QueryRecipes finds a recipe (and possible alternatives) based on the given query.
	QueryRecipes(query string) (*models.RecipeQueryResponse, error)
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
// @Description Queries recipes matching the search criteria.
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

	resp, err := h.service.QueryRecipes(payload.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
