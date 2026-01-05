// Package routes defines HTTP route registration for the application.
package routes

import (
	"golang-connect-marketplace/internal/auth/dto"
	"golang-connect-marketplace/internal/auth/http/handlers"
	"golang-connect-marketplace/internal/auth/middleware"
	"golang-connect-marketplace/internal/auth/service"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers authentication-related HTTP routes.
func RegisterRoutes(e *echo.Echo, h *handlers.Handler, authSvc *service.Service) {
	auth := e.Group("api/v1/auth")

	auth.POST("/register", h.HandleRegister)
	auth.POST("/login", h.HandleLogin)
	auth.POST("/refresh", h.HandleRefresh)
	auth.POST("/logout", h.HandleLogout)

	auth.GET("/secret", h.HandleSecret, middleware.AuthenticateMiddleware(authSvc))
	auth.GET(
		"/secret-admin",
		h.HandleSecret,
		middleware.AuthenticateMiddleware(authSvc, dto.UserRoleAdmin),
	)

	auth.GET("/github", h.HandleGithub)
	auth.GET("/github/callback", h.HandleGithubCallback)

	auth.GET("/google", h.HandleGoogle)
	auth.GET("/google/callback", h.HandleGoogleCallback)
}
