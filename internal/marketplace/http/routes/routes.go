// Package routes defines HTTP route registration for the application.
package routes

import (
	"golang-connect-marketplace/internal/auth/dto"
	"golang-connect-marketplace/internal/auth/middleware"
	"golang-connect-marketplace/internal/auth/service"
	"golang-connect-marketplace/internal/marketplace/http/handlers"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers marketplace-related HTTP routes.
func RegisterRoutes(e *echo.Echo, h *handlers.ListingsHandler, authSvc *service.Service) {
	listings := e.Group("api/v1/listings")
	cats := e.Group("api/v1/categories")

	listings.GET("", h.HandleGetListings)
	listings.POST("", h.HandleCreateListing, middleware.AuthenticateMiddleware(authSvc))
	listings.POST(
		"/:listing_id/images",
		h.HandleAddImages,
		middleware.AuthenticateMiddleware(authSvc),
	)
	listings.GET("/:listing_id", h.HandleGetListing)

	cats.POST(
		"",
		h.HandleCreateCategory,
		middleware.AuthenticateMiddleware(authSvc, dto.UserRoleAdmin),
	)
	cats.GET("", h.HandleGetCategories)
}
