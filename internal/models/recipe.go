package models

import "time"

// Recipe represents a cooking recipe in the system.
// swagger:model Recipe
type Recipe struct {
	ID                string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Title             string `gorm:"not null"`
	Ingredients       string `gorm:"type:text"`
	Steps             string `gorm:"type:text"`
	NutritionalInfo   string `gorm:"type:text"`
	AllergyDisclaimer string `gorm:"type:text"`
	Appliances        string `gorm:"type:text"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	// Additional fields can be added here (e.g., user preferences, history, pantry tracking placeholders)
}

// RecipeQueryResponse is a container for query results.
// swagger:model RecipeQueryResponse
type RecipeQueryResponse struct {
	Recipes []*Recipe `json:"recipes"`
}
