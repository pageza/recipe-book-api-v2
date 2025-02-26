package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/pageza/recipe-book-api-v2/docs"
	"github.com/pageza/recipe-book-api-v2/internal/config"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
	"github.com/pageza/recipe-book-api-v2/internal/middleware"
	"github.com/pageza/recipe-book-api-v2/internal/routes/protectedroutes"
	"github.com/pageza/recipe-book-api-v2/internal/routes/publicroutes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter initializes the Gin engine, sets up routes, and returns the router.
func NewRouter(cfg *config.Config, h *handlers.Handlers) *gin.Engine {
	// Create a new Gin engine.
	router := gin.Default()

	// Apply global middleware.
	router.Use(middleware.Logger())

	// Healthcheck endpoint.
	// @Summary Healthcheck
	// @Description Returns OK if the API is running.
	// @Tags Health
	// @Produce plain
	// @Success 200 {string} string "OK"
	// @Router /healthcheck [get]
	router.GET("/healthcheck", func(c *gin.Context) {
		log.Println("[Healthcheck] Received request for /healthcheck from", c.ClientIP())
		// Optionally, log headers.
		for key, values := range c.Request.Header {
			for _, value := range values {
				log.Printf("[Healthcheck] Header: %s = %s", key, value)
			}
		}
		c.String(http.StatusOK, "OK")
		log.Println("[Healthcheck] Responded with OK")
	})

	// Delegate public route registration to publicroutes package.
	publicroutes.Register(router, h)
	// Delegate protected route registration to protectedroutes package.
	protectedroutes.Register(router, cfg, h)
	// Serve the Swagger UI at /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Router routes registered")
	return router
}
