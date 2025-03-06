// File: internal/service/recipe_service.go
package service

import (
	"fmt"

	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
)

// RecipeService defines the interface for recipe operations.
type RecipeService interface {
	// GetRecipe retrieves a recipe by its ID.
	GetRecipe(recipeID string) (*models.Recipe, error)
	// QueryRecipes processes query requests and returns matching recipes.
	QueryRecipes(req *models.RecipeQueryRequest) (*models.RecipeQueryResponse, error)
}

// recipeService implements RecipeService.
type recipeService struct {
	repo repository.RecipeRepository
}

// NewRecipeService creates a new RecipeService instance.
func NewRecipeService(repo repository.RecipeRepository) RecipeService {
	return &recipeService{repo: repo}
}

// GetRecipe retrieves a recipe by its ID via the repository.
func (s *recipeService) GetRecipe(recipeID string) (*models.Recipe, error) {
	return s.repo.GetRecipeByID(recipeID)
}

// QueryRecipes processes the unified query request by delegating to the repository.
func (s *recipeService) QueryRecipes(req *models.RecipeQueryRequest) (*models.RecipeQueryResponse, error) {
	recipes, total, err := s.repo.QueryRecipes(req.Query, req.UserID, req.Filter, req.Page, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("repository query error: %v", err)
	}
	return &models.RecipeQueryResponse{
		Recipes: recipes,
		Page:    req.Page,
		Limit:   req.Limit,
		Total:   total,
	}, nil
}

// // Dummy helper to generate an ID.
// func generateID() string {
// 	// TODO: Use a proper UUID generator in production code.
// 	return "some-generated-id"
// }
