package service_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	"github.com/pageza/recipe-book-api-v2/pkg/utils" // gRPC client import
)

// inMemoryUserRepo is an in‑memory implementation of the UserRepository interface.
type inMemoryUserRepo struct {
	mu    sync.Mutex
	users map[string]*models.User
}

func newInMemoryUserRepo() *inMemoryUserRepo {
	return &inMemoryUserRepo{
		users: make(map[string]*models.User),
	}
}

func (r *inMemoryUserRepo) CreateUser(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.ID]; exists {
		return errors.New("user already exists")
	}
	r.users[user.ID] = user
	return nil
}

func (r *inMemoryUserRepo) GetUserByID(id string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *inMemoryUserRepo) DeleteUser(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil
}

// GetUserByEmail iterates over the stored users and returns the one with the matching email.
func (r *inMemoryUserRepo) GetUserByEmail(email string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// UpdateUser updates an existing user in the repository.
func (r *inMemoryUserRepo) UpdateUser(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	r.users[user.ID] = user
	return nil
}

// Implement other methods such as GetUserByEmail, UpdateUser as needed.

func TestUserService_Register(t *testing.T) {
	// Set up the fake repository
	repo := &inMemoryUserRepo{}
	svc := service.NewUserService(repo)

	// Create a user with hashed password.
	plainPassword := "testpassword"
	hash, err := utils.HashPassword(plainPassword)
	assert.NoError(t, err)

	user := &models.User{
		ID:           uuid.New().String(), // Use UUID for unique ID
		Username:     "testuser",
		Email:        "testuser@example.com",
		PasswordHash: hash,
		Preferences:  "{\"diet\":\"vegan\"}",
	}

	// Registration should succeed the first time.
	err = svc.Register(user)
	assert.NoError(t, err)

	// Trying to register the same user again should fail.
	err = svc.Register(user)
	assert.Error(t, err)
}

func TestUserService_Login(t *testing.T) {
	// Set up the fake repository
	repo := &inMemoryUserRepo{}
	svc := service.NewUserService(repo)

	plainPassword := "testpassword"
	hash, err := utils.HashPassword(plainPassword)
	assert.NoError(t, err)

	user := &models.User{
		ID:           uuid.New().String(),
		Username:     "testuser",
		Email:        "testuser@example.com",
		PasswordHash: hash,
		Preferences:  "{\"diet\":\"vegan\"}",
	}
	// Pre-register the user in the fake repository.
	err = repo.CreateUser(user)
	assert.NoError(t, err)

	// Login with correct credentials.
	loggedInUser, err := svc.Login("testuser@example.com", plainPassword)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, loggedInUser.Email)

	// Login with wrong password.
	_, err = svc.Login("testuser@example.com", "wrongpassword")
	assert.Error(t, err)
}

func TestUserService_UpdateAndDelete(t *testing.T) {
	// Set up the fake repository
	repo := &inMemoryUserRepo{}
	svc := service.NewUserService(repo)

	// Create a test user
	user := &models.User{
		ID:           uuid.New().String(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Preferences:  "{\"diet\":\"vegan\"}",
	}

	// Register the user first
	err := svc.Register(user)
	assert.NoError(t, err)

	// Update the user
	updatedUser := &models.User{
		ID:           user.ID,
		Username:     "updateduser",
		Email:        "updated@example.com",
		PasswordHash: "newhashedpassword",
		Preferences:  "{\"diet\":\"vegetarian\"}",
	}
	err = svc.UpdateUser(updatedUser)
	assert.NoError(t, err)

	// Verify the update
	updated, err := repo.GetUserByEmail("updated@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", updated.Username)

	// Delete the user
	err = svc.DeleteUser(user.ID)
	assert.NoError(t, err)

	// Verify the deletion
	_, err = repo.GetUserByEmail("updated@example.com")
	assert.Error(t, err, "Expected error when fetching deleted user")
}
