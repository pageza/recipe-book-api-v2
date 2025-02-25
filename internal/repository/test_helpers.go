// internal/repository/test_helpers.go
package repository

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB = *gorm.DB

// ConnectTestDB connects to the test database with retries.
func ConnectTestDB() (*gorm.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost" // default for testing
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "recipe_db"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		// If Docker Compose sets DB_PASSWORD, it gets used.
		// Otherwise fallback to "postgres" or some local default.
		dbPassword = "postgres"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var db *gorm.DB
	var err error
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		time.Sleep(3 * time.Second)
	}
	return nil, err
}
