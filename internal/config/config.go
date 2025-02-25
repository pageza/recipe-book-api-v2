/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/
package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds app configuration
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
	// Print out the config to confirm values at runtime.
	log.Printf("Connecting to DB with config: %+v\n", cfg)

	// Building the DSN from config fields. If you want to use DatabaseURL directly, swap this out.
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC options='-c search_path=public'",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)
	log.Printf("Using DSN: %s\n", dsn)

	// Very verbose GORM logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info, // Log all SQL queries at the Info level
			Colorful: false,       // If logs get messy in CI, disabling color can help
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("error connecting to DB: %w", err)
	}

	// Returning the DB without migrations here. If your migrations happen elsewhere, that's fine.
	// Or call db.AutoMigrate(...) or run your migration logic if desired.
	log.Println("Connected to the database successfully!")
	return db, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
