package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/pageza/recipe-book-api-v2/docs"
	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	"github.com/pageza/recipe-book-api-v2/internal/middleware"
	"github.com/pageza/recipe-book-api-v2/internal/routes/protectedroutes"
	"github.com/pageza/recipe-book-api-v2/internal/routes/publicroutes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// NewRouter initializes the Gin engine, sets up routes, and returns the router.
func NewRouter(cfg *config.Config, h *handlers.Handlers) *gin.Engine {
	// Create a new Gin engine.
	router := gin.Default()

	// Apply global middleware.
	router.Use(middleware.Logger())

	// Healthcheck endpoint using the dedicated handler.
	router.GET("/healthcheck", handlers.HealthHandler)

	// Delegate public route registration to publicroutes package.
	publicroutes.Register(router, h)
	// Delegate protected route registration to protectedroutes package.
	protectedroutes.Register(router, cfg, h)
	// Serve the Swagger UI at /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	zap.L().Info("Router routes registered")
	return router
}
