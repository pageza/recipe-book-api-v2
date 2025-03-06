package models

import "time"

// NutritionalInfo represents nutritional information for a recipe.
type NutritionalInfo struct {
	Calories      float64 `json:"calories"`
	Protein       float64 `json:"protein"`
	Carbohydrates float64 `json:"carbohydrates"`
	Fat           float64 `json:"fat"`
	Fiber         float64 `json:"fiber"`
}

// Recipe represents the domain model for a recipe.
type Recipe struct {
	ID                string          `json:"id" gorm:"primaryKey"`
	Title             string          `json:"title"`
	Ingredients       []string        `json:"ingredients" gorm:"type:text[]"`
	Steps             []string        `json:"steps" gorm:"type:text[]"`
	NutritionalInfo   NutritionalInfo `json:"nutritional_info" gorm:"embedded;embeddedPrefix:nutri_"`
	AllergyDisclaimer string          `json:"allergy_disclaimer"`
	Appliances        []string        `json:"appliances" gorm:"type:text[]"`
	CreatedAt         time.Time       `json:"created_at"` // time of creation
	UpdatedAt         time.Time       `json:"updated_at"` // time of last update
	UserID            string          `json:"user_id,omitempty"`
}

// RecipeQueryRequest carries parameters for querying recipes.
// An empty Query denotes a simple listing, while a non-empty value
// triggers advanced search logic.
type RecipeQueryRequest struct {
	Query  string `json:"query"`
	UserID string `json:"user_id,omitempty"`
	Filter string `json:"filter,omitempty"`
	Page   int    `json:"page,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

// RecipeQueryResponse represents the response structure for recipe queries.
type RecipeQueryResponse struct {
	Recipes []*Recipe `json:"recipes"`
	Page    int       `json:"page"`  // current page number
	Limit   int       `json:"limit"` // number of recipes per page
	Total   int       `json:"total"` // total recipes matching the query
}
