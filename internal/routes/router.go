package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	"github.com/pageza/recipe-book-api-v2/internal/middleware"
	"github.com/pageza/recipe-book-api-v2/internal/routes/protectedroutes"
	"github.com/pageza/recipe-book-api-v2/internal/routes/publicroutes"
)

// NewRouter initializes the Gin engine, sets up routes, and returns the router.
// NewRouter initializes the Gin engine, sets up routes, and returns the router.
func NewRouter(cfg *config.Config, h *handlers.Handlers) *gin.Engine {
	// Create a new Gin engine
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.Logger())

	// Healthcheck endpoint (could be a simple Gin handler as well)
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Delegate public route registration to publicroutes package
	publicroutes.Register(router, h)
	// Delegate protected route registration to protectedroutes package
	protectedroutes.Register(router, cfg, h)

	return router
}
