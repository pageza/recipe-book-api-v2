package service_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	"github.com/pageza/recipe-book-api-v2/pkg/utils" // gRPC client import
)

// fakeUserRepository implements service.UserRepository for testing.
type fakeUserRepository struct {
	users map[string]*models.User
}

func (f *fakeUserRepository) CreateUser(user *models.User) error {
	if f.users == nil {
		f.users = make(map[string]*models.User)
	}
	if _, exists := f.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	f.users[user.Email] = user
	return nil
}

func (f *fakeUserRepository) GetUserByEmail(email string) (*models.User, error) {
	if user, exists := f.users[email]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (f *fakeUserRepository) GetUserByID(id string) (*models.User, error) {
	for _, u := range f.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}
func (f *fakeUserRepository) UpdateUser(user *models.User) error {
	if f.users == nil {
		return errors.New("user not found")
	}
	// Find the user by ID first
	for _, u := range f.users {
		if u.ID == user.ID {
			// Update the user in the map using email as key
			f.users[user.Email] = user
			return nil
		}
	}
	return errors.New("user not found")
}

func (f *fakeUserRepository) DeleteUser(userID string) error {
	if f.users == nil {
		return errors.New("user not found")
	}
	// Find and delete the user by ID
	for email, u := range f.users {
		if u.ID == userID {
			delete(f.users, email)
			return nil
		}
	}
	return errors.New("user not found")
}
func TestUserService_Register(t *testing.T) {
	// Set up the fake repository
	repo := &fakeUserRepository{}
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
	repo := &fakeUserRepository{}
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
	repo := &fakeUserRepository{
		users: make(map[string]*models.User),
	}
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
