package repository

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase initializes a PostgreSQL database connection.
func NewDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}
	return db, nil
}
