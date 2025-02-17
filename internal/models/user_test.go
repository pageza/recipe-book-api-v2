package models_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUserModelInitialization(t *testing.T) {
	// Create a new user with generated ID
	id := uuid.New().String()
	user := models.User{
		ID:           id,
		Username:     "testuser",
		Email:        "testuser@example.com",
		PasswordHash: "$2a$14$abcdefghijklmnopqrstuv",
		Preferences:  "{\"diet\":\"vegan\"}",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Verify that fields are assigned correctly
	assert.Equal(t, id, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "testuser@example.com", user.Email)
	assert.Contains(t, user.PasswordHash, "$2a$14$")
	assert.Equal(t, "{\"diet\":\"vegan\"}", user.Preferences)
	// Now we assert that the timestamps are not zero.
	assert.False(t, user.CreatedAt.IsZero(), "CreatedAt should be set")
	assert.False(t, user.UpdatedAt.IsZero(), "UpdatedAt should be set")
}
