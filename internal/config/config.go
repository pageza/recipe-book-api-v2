/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/
// config/config.go
package config

import (
	"fmt"
	"os"

	"github.com/pageza/recipe-book-api-v2/internal/models" // ✅ Import Notification model
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DatabaseURL string
	Port        string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        getEnv("PORT", "8080"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "recipe_db"),
		JWTSecret:   getEnv("JWT_SECRET", "your_jwt_secret"),
	}
	return cfg, nil
}

func ConnectDatabase(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// ✅ Auto-migrate the notifications table
	err = db.AutoMigrate(&models.Notification{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
