// File: internal/service/recipe_service.go
package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
)

// RecipeService defines the interface for recipe operations.
type RecipeService interface {
	CreateRecipe(recipe *models.Recipe) error
	// GetRecipe retrieves a recipe by its ID.
	GetRecipe(recipeID string) (*models.Recipe, error)
	// ListRecipes returns all stored recipes.
	ListRecipes() ([]*models.Recipe, error)
	// GetAllRecipes is an alias to ListRecipes.
	GetAllRecipes() ([]*models.Recipe, error)
	// GetRecipeByQuery returns a recipe based on a query string.
	GetRecipeByQuery(query string) (*models.Recipe, error)
}

// recipeService implements RecipeService.
type recipeService struct {
	repo repository.RecipeRepository
}

// NewRecipeService creates a new RecipeService instance.
func NewRecipeService(repo repository.RecipeRepository) RecipeService {
	return &recipeService{repo: repo}
}

// CreateRecipe adds a new recipe to the database.

func (s *recipeService) CreateRecipe(recipe *models.Recipe) error {
	if recipe.Title == "" {
		return errors.New("recipe title cannot be empty")
	}
	if recipe.ID == "" {
		recipe.ID = uuid.New().String()
	}
	return s.repo.CreateRecipe(recipe)
}

// GetRecipe retrieves a recipe by ID.
func (s *recipeService) GetRecipe(recipeID string) (*models.Recipe, error) {
	return s.repo.GetRecipeByID(recipeID)
}

// ListRecipes returns all stored recipes.
func (s *recipeService) ListRecipes() ([]*models.Recipe, error) {
	return s.repo.GetAllRecipes()
}

// GetAllRecipes returns all stored recipes (alias to ListRecipes).
func (s *recipeService) GetAllRecipes() ([]*models.Recipe, error) {
	return s.ListRecipes()
}

// GetRecipeByQuery returns a recipe based on the query string.
// For now, this is a stub that returns an error to indicate it is not implemented.
func (s *recipeService) GetRecipeByQuery(query string) (*models.Recipe, error) {
	// TODO: Implement query logic using your RAG system.
	return nil, errors.New("query recipe not implemented")
}

// Dummy helper to generate an ID.
func generateID() string {
	// TODO: Use a proper UUID generator in production code.
	return "some-generated-id"
}
