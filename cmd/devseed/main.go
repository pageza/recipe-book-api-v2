package main

import (
	"log"
	"time"

	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
)

func main() {
	// Load configuration and connect to the database.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the Recipe model to ensure the table exists.
	if err := db.AutoMigrate(&models.Recipe{}); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	// Check if there are any recipes already.
	var count int64
	db.Model(&models.Recipe{}).Count(&count)
	if count == 0 {
		log.Println("No recipes found. Inserting dummy recipes for development.")
		dummyRecipes := getDummyRecipes()
		repo := repository.NewRecipeRepository(db)
		for _, recipe := range dummyRecipes {
			// Set initial timestamps.
			recipe.CreatedAt = time.Now()
			recipe.UpdatedAt = time.Now()
			if err := repo.CreateRecipe(&recipe); err != nil {
				log.Printf("Failed to insert recipe: %s, error: %v", recipe.Title, err)
			}
		}
		log.Println("Dummy recipes inserted.")
	} else {
		log.Println("Recipes already exist, skipping dummy seeding.")
	}
}

func getDummyRecipes() []models.Recipe {
	return []models.Recipe{
		{
			Title: "Spaghetti Bolognese",
			Ingredients: `{"items": [
				{"name": "spaghetti", "quantity": "200g"},
				{"name": "ground beef", "quantity": "300g"},
				{"name": "tomato sauce", "quantity": "500ml"}
			]}`,
			Steps: `{"steps": [
				"Boil water and cook spaghetti until al dente",
				"Cook ground beef with onions and garlic in a pan",
				"Mix beef with tomato sauce",
				"Combine spaghetti with sauce and serve hot"
			]}`,
			NutritionalInfo:   `{"calories": "700", "protein": "30g", "fat": "20g"}`,
			AllergyDisclaimer: "Contains gluten and dairy",
			Appliances:        `["Stove", "Pot", "Pan"]`,
		},
		{
			Title: "Vegetarian Stir-fry",
			Ingredients: `{"items": [
				{"name": "tofu", "quantity": "200g"},
				{"name": "broccoli", "quantity": "150g"},
				{"name": "bell pepper", "quantity": "1 piece"}
			]}`,
			Steps: `{"steps": [
				"Press tofu to remove excess water",
				"Stir-fry tofu until golden in a wok",
				"Add chopped broccoli and bell pepper",
				"Season with soy sauce and spices, and cook until vegetables are tender"
			]}`,
			NutritionalInfo:   `{"calories": "400", "protein": "15g", "fat": "10g"}`,
			AllergyDisclaimer: "Contains soy",
			Appliances:        `["Wok", "Spatula"]`,
		},
	}
}
