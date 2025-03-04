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

// ResolveRecipeQuery orchestrates the full recipe retrieval flow.
func (s *recipeService) ResolveRecipeQuery(query string) (*models.RecipeQueryResponse, error) {
	// 1. Retrieve all recipes from PostgreSQL.
	recipes, err := s.ListRecipes()
	if err != nil {
		return nil, err
	}

	// 2. Perform a simple substring match to find a recipe.
	var matchedRecipe *models.Recipe
	qLower := strings.ToLower(query)
	for _, r := range recipes {
		if strings.Contains(strings.ToLower(r.Title), qLower) ||
			strings.Contains(strings.ToLower(r.Ingredients), qLower) ||
			strings.Contains(strings.ToLower(r.Steps), qLower) {
			matchedRecipe = r
			break
		}
	}

	// 3a. If a matching recipe is found, return it.
	if matchedRecipe != nil {
		return &models.RecipeQueryResponse{
			Recipes: []*models.Recipe{matchedRecipe},
		}, nil
	}

	// 3b. If no match is found, generate a new recipe via AI (stubbed).
	newRecipe, err := s.generateRecipeFromQuery(query)
	if err != nil {
		return nil, err
	}

	// 4. Persist the newly generated recipe in PostgreSQL.
	if err := s.CreateRecipe(newRecipe); err != nil {
		return nil, err
	}

	// 5. Update the vector embedding (placeholder for future PGVector integration).
	if err := s.updateRecipeVectorEmbedding(newRecipe); err != nil {
		zap.L().Warn("Vector DB update failed", zap.String("recipeID", newRecipe.ID), zap.Error(err))
	}

	// 6. Return the newly generated recipe.
	return &models.RecipeQueryResponse{
		Recipes: []*models.Recipe{newRecipe},
	}, nil
}

// generateRecipeFromQuery is a stub simulating an AI service call to create a new recipe.
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

// updateRecipeVectorEmbedding is a placeholder for future vector DB integration.
func (s *recipeService) updateRecipeVectorEmbedding(recipe *models.Recipe) error {
	zap.L().Info("Vector embedding update for recipe (stub)", zap.String("recipeID", recipe.ID))
	return nil
}

// Dummy helper to generate an ID.
func generateID() string {
	// TODO: Use a proper UUID generator in production code.
	return "some-generated-id"
}
