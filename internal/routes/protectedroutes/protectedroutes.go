package protectedroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	"github.com/pageza/recipe-book-api-v2/internal/handlers/recipes"
	"github.com/pageza/recipe-book-api-v2/internal/middleware"
)

// Register registers protected endpoints.
// Since the resolver handles recipe generation via the query endpoint,
// we only expose endpoints to query recipes and retrieve stored recipes.
func Register(router *gin.Engine, cfg *config.Config, h *handlers.Handlers, recipeHandler *recipes.RecipeHandler) {
	protected := router.Group("/")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		// User endpoint.
		protected.GET("/profile", h.User.Profile)

		// Recipe endpoints.
		// The user submits a query that is processed by the resolver logic.
		protected.POST("/recipe/query", h.Recipe.Query)
		// Retrieve a specific recipe by its ID.
		protected.GET("/recipe/:id", h.Recipe.Get)
		// List all recipes (e.g., those previously generated for the logged-in user).
		protected.GET("/recipes", recipeHandler.Query)
	}
}
