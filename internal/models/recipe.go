package models

// Recipe represents a recipe in the database.
type Recipe struct {
	ID          string `gorm:"primaryKey"`
	Title       string
	Ingredients string // JSON string of ingredient list
	Steps       string // JSON string of step list
}
