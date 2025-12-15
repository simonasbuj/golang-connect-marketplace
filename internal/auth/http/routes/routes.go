package routes

import (
    "github.com/labstack/echo/v4"
    "golang-connect-marketplace/internal/auth/http/handlers"
)

func RegisterRoutes(e *echo.Echo, h *handlers.AuthHandler) {
	auth := e.Group("api/v1/auth")

    auth.POST("/register", h.HandleRegister)
    auth.POST("/login", h.HandleLogin)
    auth.POST("/refresh", h.HandleRefresh)
}
