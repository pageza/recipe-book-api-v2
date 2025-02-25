package publicroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/pageza/recipe-book-api-v2/internal/handlers"
)

// Register registers public routes and accepts the composite handlers.
func Register(router *gin.Engine, h *handlers.Handlers) {
	router.POST("/register", h.User.Register)
	router.POST("/login", h.User.Login)
	router.POST("/request-password-reset", h.User.RequestPasswordReset)
	router.POST("/reset-password", h.User.ResetPassword)
}
