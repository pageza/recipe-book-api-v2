package protectedroutes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	"github.com/pageza/recipe-book-api-v2/internal/middleware"
)

// Register registers protected routes and accepts both config and handlers.
func Register(router *gin.Engine, cfg *config.Config, h *handlers.Handlers) {
	protected := router.Group("/")
	protected.Use(
		middleware.JWTAuth(cfg.JWTSecret),
		middleware.RateLimitMiddleware(1*time.Second, 5), // For example, allow 5 requests per second.
		middleware.AuditMiddleware(),
	)
	{
		protected.GET("/profile", h.User.Profile)
		// Add more protected routes here.
	}
}
