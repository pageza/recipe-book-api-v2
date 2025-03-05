package protectedroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	"github.com/pageza/recipe-book-api-v2/internal/middleware"
)

// Register registers protected routes and accepts both config and handlers.
func Register(router *gin.Engine, cfg *config.Config, h *handlers.Handlers) {
	// Group protected routes with JWT authentication middleware.
	protected := router.Group("/")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		protected.GET("/profile", h.User.Profile)

		// Register recipe endpoints under the protected routes.
		recipes := protected.Group("api/v1/recipes")
		{
			recipes.POST("/query", h.Recipe.Query)
			recipes.POST("/", h.Recipe.Create)
			recipes.GET("/:id", h.Recipe.Get)
			recipes.GET("/", h.Recipe.List)
		}

		// Add additional protected routes as needed.
	}
}
