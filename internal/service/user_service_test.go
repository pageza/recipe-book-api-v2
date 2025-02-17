package service_test

import (
	"errors"
	"testing"

	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	"github.com/pageza/recipe-book-api-v2/pkg/utils"
	"github.com/stretchr/testify/assert"
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

func TestUserService_Register(t *testing.T) {
	repo := &fakeUserRepository{}
	svc := service.NewUserService(repo)

	// Create a user with hashed password.
	plainPassword := "testpassword"
	hash, err := utils.HashPassword(plainPassword)
	assert.NoError(t, err)

	user := &models.User{
		ID:           "user-1",
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
	repo := &fakeUserRepository{}
	svc := service.NewUserService(repo)

	plainPassword := "testpassword"
	hash, err := utils.HashPassword(plainPassword)
	assert.NoError(t, err)

	user := &models.User{
		ID:           "user-1",
		Username:     "testuser",
		Email:        "testuser@example.com",
		PasswordHash: hash,
		Preferences:  "{\"diet\":\"vegan\"}",
	}
	// Pre-register the user.
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
