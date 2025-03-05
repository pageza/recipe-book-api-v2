// File: internal/service/recipe_service.go
package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
	"go.uber.org/zap"
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
	// QueryRecipes returns a list of recipes matching the query.
	QueryRecipes(query string) (*models.RecipeQueryResponse, error)
	ResolveRecipeQuery(query string) (*models.RecipeQueryResponse, error)
}

// recipeService implements RecipeService.
type recipeService struct {
	repo repository.RecipeRepository
}

// NewRecipeService creates a new RecipeService instance.
func NewRecipeService(repo repository.RecipeRepository) RecipeService {
	return &recipeService{repo: repo}
}

// CreateRecipe creates a new recipe in the database.
func (s *recipeService) CreateRecipe(recipe *models.Recipe) error {
	if recipe.Title == "" {
		return errors.New("recipe title cannot be empty")
	}
	if recipe.ID == "" {
		recipe.ID = uuid.New().String()
	}
	if err := s.repo.CreateRecipe(recipe); err != nil {
		zap.L().Error("failed to create recipe", zap.Error(err))
		return err
	}
	if err := s.updateRecipeVectorEmbedding(recipe); err != nil {
		zap.L().Warn("failed to update recipe vector embedding", zap.Error(err))
	}
	return nil
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
// Stub implementation; full query logic to be implemented later.
func (s *recipeService) GetRecipeByQuery(query string) (*models.Recipe, error) {
	return nil, errors.New("query recipe not implemented")
}

// QueryRecipes returns a response containing recipes matching the query.
func (s *recipeService) QueryRecipes(query string) (*models.RecipeQueryResponse, error) {
	recipes, err := s.ListRecipes()
	if err != nil {
		return nil, err
	}
	return &models.RecipeQueryResponse{Recipes: recipes}, nil
}

// ResolveRecipeQuery handles a query and returns a matching recipe or generates one.
func (s *recipeService) ResolveRecipeQuery(query string) (*models.RecipeQueryResponse, error) {
	recipes, err := s.repo.GetAllRecipes()
	if err != nil {
		return nil, err
	}

	var matched *models.Recipe
	lowerQuery := strings.ToLower(query)
	for _, r := range recipes {
		if strings.Contains(strings.ToLower(r.Title), lowerQuery) ||
			strings.Contains(strings.ToLower(r.Ingredients), lowerQuery) ||
			strings.Contains(strings.ToLower(r.Steps), lowerQuery) {
			matched = r
			break
		}
	}

	if matched != nil {
		return &models.RecipeQueryResponse{
			Recipes: []*models.Recipe{matched},
		}, nil
	}

	generated, err := s.generateRecipeFromQuery(query)
	if err != nil {
		return nil, err
	}

	if err := s.CreateRecipe(generated); err != nil {
		return nil, err
	}

	if err := s.updateRecipeVectorEmbedding(generated); err != nil {
		zap.L().Warn("failed to update vector embedding for generated recipe", zap.Error(err))
	}

	return &models.RecipeQueryResponse{
		Recipes: []*models.Recipe{generated},
	}, nil
}

// generateRecipeFromQuery creates a dummy recipe based on the query string.
func (s *recipeService) generateRecipeFromQuery(query string) (*models.Recipe, error) {
	newID := uuid.New().String()
	return &models.Recipe{
		ID:                newID,
		Title:             fmt.Sprintf("%s - Generated Recipe", query),
		Ingredients:       "[]",
		Steps:             "[]",
		NutritionalInfo:   "{}",
		AllergyDisclaimer: "",
		Appliances:        "[]",
	}, nil
}

// updateRecipeVectorEmbedding is a stub for vector embedding update logic.
func (s *recipeService) updateRecipeVectorEmbedding(recipe *models.Recipe) error {
	zap.L().Info("Updating vector embedding for recipe", zap.String("id", recipe.ID))
	return nil
}

// Dummy helper to generate an ID.
func generateID() string {
	// TODO: Use a proper UUID generator in production code.
	return "some-generated-id"
}
