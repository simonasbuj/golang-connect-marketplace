// Package routes defines HTTP route registration for the application.
package routes

import (
	"golang-connect-marketplace/internal/auth/middleware"
	"golang-connect-marketplace/internal/auth/service"
	"golang-connect-marketplace/internal/marketplace/http/handlers"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers marketplace-related HTTP routes.
func RegisterRoutes(e *echo.Echo, h *handlers.ListingsHandler, authSvc *service.Service) {
	listingsAPI := e.Group("api/v1/listings")

	listingsAPI.GET("", h.HandleGetListings)
	listingsAPI.POST("", h.HandleCreateListing, middleware.AuthenticateMiddleware(authSvc))
}
