package repository

import (
	"fmt"

	"github.com/pageza/recipe-book-api-v2/internal/models"
	"gorm.io/gorm"
)

// RecipeRepository defines the data access interface for recipes.
type RecipeRepository interface {
	// GetRecipeByID retrieves a recipe by its unique ID.
	GetRecipeByID(recipeID string) (*models.Recipe, error)
	// QueryRecipes performs a search and filtering query on recipes.
	// It returns the matched recipes, the total count for pagination, and an error (if any).
	QueryRecipes(query, userID, filter string, page, limit int) ([]*models.Recipe, int, error)
}

// recipeRepository is the struct that implements RecipeRepository
type recipeRepository struct {
	db *gorm.DB
}

// NewRecipeRepository returns an implementation of RecipeRepository
func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &recipeRepository{db: db}
}

// GetRecipeByID retrieves a recipe by its ID.
func (r *recipeRepository) GetRecipeByID(recipeID string) (*models.Recipe, error) {
	var recipe models.Recipe
	if err := r.db.First(&recipe, "id = ?", recipeID).Error; err != nil {
		return nil, err
	}
	return &recipe, nil
}

// QueryRecipes performs a query with optional filters:
//   - If userID is provided, it filters by recipe creator.
//   - If filter is provided, it applies additional filtering on the title.
//   - If query text is provided, it searches in title and ingredients.
//
// Pagination is applied via page and limit parameters.
func (r *recipeRepository) QueryRecipes(query, userID, filter string, page, limit int) ([]*models.Recipe, int, error) {
	var recipes []*models.Recipe
	dbQuery := r.db.Model(&models.Recipe{})

	if userID != "" {
		dbQuery = dbQuery.Where("user_id = ?", userID)
	}
	if filter != "" {
		dbQuery = dbQuery.Where("title LIKE ?", "%"+filter+"%")
	}
	if query != "" {
		// For advanced search, search in both title and ingredients.
		dbQuery = dbQuery.Where("title LIKE ? OR ingredients LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	// Retrieve the total count before pagination.
	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count recipes: %v", err)
	}

	// Calculate the offset based on the page number.
	offset := (page - 1) * limit
	if err := dbQuery.Offset(offset).Limit(limit).Find(&recipes).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to query recipes: %v", err)
	}

	return recipes, int(total), nil
}
