// Package routes defines HTTP route registration for the application.
package routes

import (
	"golang-connect-marketplace/internal/auth/dto"
	m "golang-connect-marketplace/internal/auth/middleware"
	"golang-connect-marketplace/internal/auth/service"
	"golang-connect-marketplace/internal/marketplace/http/handlers"

	"github.com/labstack/echo/v4"
)

// RegisterListingsRoutes registers marketplace-related HTTP routes.
func RegisterListingsRoutes(e *echo.Echo, lh *handlers.ListingsHandler, authSvc *service.Service) {
	listings := e.Group("api/v1/listings")
	cats := e.Group("api/v1/categories")

	listings.GET("", lh.HandleGetListings)
	listings.POST("", lh.HandleCreateListing, m.AuthenticateMiddleware(authSvc))
	listings.GET("/:listing_id", lh.HandleGetListing)
	listings.POST("/:listing_id/images", lh.HandleAddImages, m.AuthenticateMiddleware(authSvc))
	listings.DELETE("/:listing_id/images", lh.HandleDeleteImages, m.AuthenticateMiddleware(authSvc))

	cats.POST("", lh.HandleCreateCategory, m.AuthenticateMiddleware(authSvc, dto.UserRoleAdmin))
	cats.GET("", lh.HandleGetCategories)
}
