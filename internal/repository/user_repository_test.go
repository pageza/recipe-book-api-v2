package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/repository"
)

func TestUserRepository_CreateAndGetUser(t *testing.T) {
	// Open an in-memory SQLite DB for testing.
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Create a new user repository (AutoMigrate will create the table).
	repo := repository.NewUserRepository(db)

	// Create a test user.
	user := &models.User{
		ID:           "dummy-uuid-1",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$14$abcdefghijklmnopqrstuv", // Not verifying password here.
		Preferences:  "{\"diet\":\"vegan\"}",
	}

	err = repo.CreateUser(user)
	assert.NoError(t, err)

	// Retrieve the user by email.
	fetched, err := repo.GetUserByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.Username, fetched.Username)
	assert.Equal(t, user.Email, fetched.Email)
	assert.Equal(t, user.PasswordHash, fetched.PasswordHash)
}
