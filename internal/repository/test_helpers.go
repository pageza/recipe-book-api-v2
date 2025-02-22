package repository

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is an alias for gorm.DB.
// In your code, repository.DB is expected to be this type.
func ConnectTestDB() (*gorm.DB, error) {
	dsn := "host=postgres user=postgres password=postgres dbname=recipe_db port=5432 sslmode=disable TimeZone=UTC"
	var db *gorm.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		// Wait for 2 seconds before retrying
		time.Sleep(2 * time.Second)
	}
	return nil, err
}
