// health.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	zap.L().Info("[Healthcheck] Received request", zap.String("client_ip", c.ClientIP()))
	// (Optional) log request headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			zap.L().Debug("[Healthcheck] Header", zap.String("key", key), zap.String("value", value))
		}
	}
	c.String(http.StatusOK, "OK")
	zap.L().Info("[Healthcheck] Responded with OK")
}
