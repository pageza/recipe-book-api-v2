package protectedroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	"github.com/pageza/recipe-book-api-v2/internal/middleware"
)

// Register registers protected routes and accepts both config and handlers.
func Register(router *gin.Engine, cfg *config.Config, h *handlers.Handlers) {
	protected := router.Group("/")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		protected.GET("/profile", h.User.Profile)
		// Add more protected routes here.
	}
}
