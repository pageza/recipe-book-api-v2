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
		PasswordHash: "$2a$14$abcdefghijklmnopqrstuv", // Simulating hashed password.
		Preferences:  "{\"diet\":\"vegan\"}",
	}

	// Create the user in the repository (using the in-memory SQLite DB).
	err = repo.CreateUser(user)
	assert.NoError(t, err, "Expected no error when creating user")

	// Retrieve the user by email from the repository.
	fetched, err := repo.GetUserByEmail("test@example.com")
	assert.NoError(t, err, "Expected no error when fetching user by email")

	// Validate that the fetched user matches the created user.
	assert.Equal(t, user.Username, fetched.Username, "Expected username to match")
	assert.Equal(t, user.Email, fetched.Email, "Expected email to match")
	assert.Equal(t, user.PasswordHash, fetched.PasswordHash, "Expected password hash to match")
}

func TestUserRepository_GetUserByEmail_NotFound(t *testing.T) {
	// Open an in-memory SQLite DB for testing.
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Create a new user repository (AutoMigrate will create the table).
	repo := repository.NewUserRepository(db)

	// Try to fetch a user that doesn't exist.
	fetched, err := repo.GetUserByEmail("nonexistent@example.com")
	assert.Error(t, err, "Expected error when fetching user by non-existent email")
	assert.Nil(t, fetched, "Expected fetched user to be nil")
}

func TestUserRepository_CreateUser_WithDuplicateEmail(t *testing.T) {
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
		PasswordHash: "$2a$14$abcdefghijklmnopqrstuv", // Simulating hashed password.
		Preferences:  "{\"diet\":\"vegan\"}",
	}

	// Create the first user.
	err = repo.CreateUser(user)
	assert.NoError(t, err, "Expected no error when creating user")

	// Attempt to create a user with the same email.
	duplicateUser := &models.User{
		ID:           "dummy-uuid-2",
		Username:     "testuser2",
		Email:        "test@example.com",              // Same email as the first user.
		PasswordHash: "$2a$14$abcdefghijklmnopqrstuv", // Simulating hashed password.
		Preferences:  "{\"diet\":\"vegetarian\"}",
	}

	// Attempt to create a user with the same email, expecting an error.
	err = repo.CreateUser(duplicateUser)
	assert.Error(t, err, "Expected error when creating user with duplicate email")
}
