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
	// Read DB_HOST from environment, default to "postgres" if not set.
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres"
	}
	dsn := fmt.Sprintf("host=%s user=postgres password=postgres dbname=recipe_db port=5432 sslmode=disable TimeZone=UTC", dbHost)
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
