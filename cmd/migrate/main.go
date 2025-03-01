// cmd/migrate/main.go
package main

import (
	"log"
	"os"

	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/models"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to the database.
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Instead of os.Getenv("CI"), check a dedicated variable:
	if os.Getenv("DROP_TABLES") == "true" {
		log.Println("DROP_TABLES environment detected, dropping existing tables")
		if err := db.Migrator().DropTable(&models.User{}, &models.Recipe{}, &models.Notification{}); err != nil {
			log.Fatalf("failed to drop tables: %v", err)
		}
	}

	// Run migrations.
	err = db.AutoMigrate(&models.User{}, &models.Recipe{}, &models.Notification{})
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Database migrations complete")
}
