// Package handlers defines HTTP handlers for application marketplace endpoints.
package handlers

import (
	"golang-connect-marketplace/internal/marketplace/services"
	"golang-connect-marketplace/pkg/responses"

	"github.com/labstack/echo/v4"
)

// ListingsHandler handles listings-related HTTP requests.
type ListingsHandler struct {
	svc *services.ListingsService
}

// NewListingsHandler creates a new Handler for handling listings requests.
func NewListingsHandler(svc *services.ListingsService) *ListingsHandler {
	return &ListingsHandler{
		svc: svc,
	}
}

// HandleCreateListing handles requests to create new listing.
func (h *ListingsHandler) HandleCreateListing(c echo.Context) error {
	return responses.JSONSuccess(c, "created new item", nil)
}

// HandleGetListings handles requests to get a list of listings.
func (h *ListingsHandler) HandleGetListings(c echo.Context) error {
	return responses.JSONSuccess(c, "fetched items", nil)
}
