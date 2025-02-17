package utils_test

import (
	"testing"

	"github.com/pageza/recipe-book-api-v2/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashing(t *testing.T) {
	plain := "testpassword"
	hash, err := utils.HashPassword(plain)
	assert.NoError(t, err, "expected no error during hashing")
	assert.NotEmpty(t, hash, "expected a non-empty hash")

	// Verify that the correct password matches
	match := utils.CheckPasswordHash(plain, hash)
	assert.True(t, match, "password should match its hash")

	// And that an incorrect password does not match
	match = utils.CheckPasswordHash("wrongpassword", hash)
	assert.False(t, match, "wrong password should not match the hash")
}
