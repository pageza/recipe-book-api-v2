package repository

import (
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"gorm.io/gorm"
)

// RecipeRepository defines the interface for database operations
type RecipeRepository interface {
	CreateRecipe(recipe *models.Recipe) error
	GetRecipeByID(recipeID string) (*models.Recipe, error)
	GetAllRecipes() ([]*models.Recipe, error)
}

// recipeRepository is the struct that implements RecipeRepository
type recipeRepository struct {
	db *gorm.DB
}

// NewRecipeRepository returns an implementation of RecipeRepository
func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &recipeRepository{db: db}
}

// Implementing interface methods
func (r *recipeRepository) CreateRecipe(recipe *models.Recipe) error {
	return r.db.Create(recipe).Error
}

func (r *recipeRepository) GetRecipeByID(recipeID string) (*models.Recipe, error) {
	var recipe models.Recipe
	err := r.db.First(&recipe, "id = ?", recipeID).Error
	return &recipe, err
}

func (r *recipeRepository) GetAllRecipes() ([]*models.Recipe, error) {
	var recipes []*models.Recipe
	err := r.db.Find(&recipes).Error
	return recipes, err
}
