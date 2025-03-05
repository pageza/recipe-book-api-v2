package repository

import (
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"gorm.io/gorm"
)

// RecipeRepository defines methods to interact with recipes in the database.
type RecipeRepository interface {
	CreateRecipe(recipe *models.Recipe) error
	GetRecipeByID(id string) (*models.Recipe, error)
	GetAllRecipes() ([]*models.Recipe, error)
}

type recipeRepository struct {
	db *gorm.DB
}

// NewRecipeRepository creates a new instance of RecipeRepository.
func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &recipeRepository{db: db}
}

// CreateRecipe inserts a new recipe record in the database.
func (r *recipeRepository) CreateRecipe(recipe *models.Recipe) error {
	return r.db.Create(recipe).Error
}

// GetRecipeByID retrieves a recipe by its ID.
func (r *recipeRepository) GetRecipeByID(id string) (*models.Recipe, error) {
	var recipe models.Recipe
	if err := r.db.First(&recipe, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &recipe, nil
}

// GetAllRecipes returns all recipes from the database.
func (r *recipeRepository) GetAllRecipes() ([]*models.Recipe, error) {
	var recipes []*models.Recipe
	if err := r.db.Find(&recipes).Error; err != nil {
		return nil, err
	}
	return recipes, nil
}
