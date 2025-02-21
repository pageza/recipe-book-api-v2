package service

import (
	"errors"

	"github.com/pageza/recipe-book-api-v2/internal/models" // ✅ Import models
	"github.com/pageza/recipe-book-api-v2/internal/repository"
)

// RecipeService defines the interface for recipe operations
type RecipeService interface {
	CreateRecipe(recipe *models.Recipe) error
	GetRecipe(recipeID string) (*models.Recipe, error)
	ListRecipes() ([]*models.Recipe, error)
}

// recipeService implements RecipeService
type recipeService struct {
	repo repository.RecipeRepository
}

// NewRecipeService creates a new RecipeService instance
func NewRecipeService(repo repository.RecipeRepository) RecipeService {
	return &recipeService{repo: repo} // ✅ Now correctly returns RecipeService interface
}

// CreateRecipe adds a new recipe to the database
func (s *recipeService) CreateRecipe(recipe *models.Recipe) error {
	if recipe.Title == "" {
		return errors.New("recipe title cannot be empty")
	}
	return s.repo.CreateRecipe(recipe)
}

// GetRecipe retrieves a recipe by ID
func (s *recipeService) GetRecipe(recipeID string) (*models.Recipe, error) {
	return s.repo.GetRecipeByID(recipeID)
}

// ListRecipes returns all stored recipes
func (s *recipeService) ListRecipes() ([]*models.Recipe, error) {
	return s.repo.GetAllRecipes()
}
