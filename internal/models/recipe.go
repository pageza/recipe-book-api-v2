package models

import "time"

// Recipe represents a comprehensive recipe in the database.
type Recipe struct {
	ID                string `gorm:"primaryKey"`
	Title             string // The title or name of the recipe.
	Ingredients       string // JSON string of the detailed ingredient list.
	Steps             string // JSON string of the step-by-step instructions.
	NutritionalInfo   string // JSON string containing nutritional breakdown (calories, macros, etc.)
	AllergyDisclaimer string // Disclaimer text about potential allergy risks.
	Appliances        string // JSON string of required appliances/cookware.
	CreatedAt         time.Time
	UpdatedAt         time.Time
	// Additional fields can be added here (e.g., user preferences, history, pantry tracking placeholders)
}

// RecipeQueryResponse is a container for query results.
type RecipeQueryResponse struct {
	Recipes []*Recipe `json:"recipes"`
}
