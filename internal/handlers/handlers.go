// internal/handlers/handlers.go
package handlers

import (
	//"github.com/pageza/recipe-book-api-v2/internal/handlers/recipes"
	"github.com/pageza/recipe-book-api-v2/internal/handlers/users"
)

type Handlers struct {
	User *users.UserHandler
	// Recipe *RecipeHandler
	// Add other handlers as needed
}
