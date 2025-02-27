// health.go
package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles the healthcheck endpoint.
//
// @Summary Healthcheck
// @Description Returns OK if the API is running.
// @Tags Health
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /healthcheck [get]
func HealthHandler(c *gin.Context) {
	log.Println("[Healthcheck] Received request for /healthcheck from:", c.ClientIP())
	// Log headers (optional)
	for key, values := range c.Request.Header {
		for _, value := range values {
			log.Printf("[Healthcheck] Header: %s = %s", key, value)
		}
	}
	c.String(http.StatusOK, "OK")
	log.Println("[Healthcheck] Responded with OK")
}
