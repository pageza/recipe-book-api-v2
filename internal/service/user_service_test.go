package service_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"
	// gRPC client import
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
	repo := newInMemoryUserRepo() // ensures repo.users is initialized
	svc := service.NewUserService(repo)
	user := &models.User{
		ID:           "test-id",
		Email:        "testuser@example.com",
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Preferences:  "{}",
	}

	err := svc.Register(user)
	assert.NoError(t, err)

	// Now, attempting to register the same user should yield an error.
	err = svc.Register(user)
	assert.Error(t, err)
}

func TestUserService_Login(t *testing.T) {
	repo := newInMemoryUserRepo()
	svc := service.NewUserService(repo)

	// Generate a valid hash for the plaintext "password".
	hashed, err := utils.HashPassword("password")
	assert.NoError(t, err)

	user := &models.User{
		ID:           "test-id",
		Email:        "testuser@example.com",
		Username:     "testuser",
		PasswordHash: hashed,
		Preferences:  "{}",
	}

	// Create the user in the repository.
	err = repo.CreateUser(user)
	assert.NoError(t, err)

	// Now perform the login with the correct password.
	loggedUser, err := svc.Login("testuser@example.com", "password")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, loggedUser.ID)
}

func TestUserService_UpdateAndDelete(t *testing.T) {
	// Use newInMemoryUserRepo to ensure repo.users is initialized.
	repo := newInMemoryUserRepo()
	svc := service.NewUserService(repo)

	// Create and register a test user.
	user := &models.User{
		ID:       "test-id",
		Email:    "test@example.com",
		Username: "testuser",
		// For Update/Delete tests we don't need a valid hashed password,
		// but if a password check was involved, use utils.HashPassword("password")
		PasswordHash: "dummy",
		Preferences:  "{}",
	}
	err := svc.Register(user)
	assert.NoError(t, err, "Register should succeed for a new user")

	// Update the user.
	user.Username = "updateduser"
	err = svc.UpdateUser(user)
	assert.NoError(t, err, "UpdateUser should succeed for an existing user")

	// Confirm the update by fetching the user from the repo.
	updatedUser, err := repo.GetUserByID(user.ID)
	assert.NoError(t, err, "GetUserByID should succeed for an updated user")
	assert.Equal(t, "updateduser", updatedUser.Username, "Username should be updated")

	// Delete the user.
	err = svc.DeleteUser(user.ID)
	assert.NoError(t, err, "DeleteUser should succeed for an existing user")

	// The deleted user should no longer be found.
	_, err = repo.GetUserByID(user.ID)
	assert.Error(t, err, "Fetching a deleted user should return an error")
}
