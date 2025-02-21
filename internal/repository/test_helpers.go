package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is an alias for gorm.DB.
// In your code, repository.DB is expected to be this type.
type DB = *gorm.DB

// ConnectTestDB connects to the test database.
// Adjust the DSN below to match your testing database configuration.
func ConnectTestDB() (*gorm.DB, error) {
	// DSN for your test database.
	// You may want to read this from an environment variable instead.
	dsn := "host=postgres user=postgres password=postgres dbname=recipe_db port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
