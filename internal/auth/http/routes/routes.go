// Package routes defines HTTP route registration for the application.
package routes

import (
	"golang-connect-marketplace/internal/auth/http/handlers"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers authentication-related HTTP routes.
func RegisterRoutes(e *echo.Echo, h *handlers.Handler) {
	auth := e.Group("api/v1/auth")

	auth.POST("/register", h.HandleRegister)
	auth.POST("/login", h.HandleLogin)
	auth.POST("/refresh", h.HandleRefresh)
}
