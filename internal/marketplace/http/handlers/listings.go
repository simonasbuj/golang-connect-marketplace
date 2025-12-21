// Package handlers defines HTTP handlers for application marketplace endpoints.
package handlers

import (
	"golang-connect-marketplace/internal/marketplace/dto"
	"golang-connect-marketplace/internal/marketplace/services"
	"golang-connect-marketplace/pkg/responses"
	"golang-connect-marketplace/pkg/validation"
	"net/http"

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

// HandleCreateCategory handles requests to create new category.
func (h *ListingsHandler) HandleCreateCategory(c echo.Context) error {
	var reqDto dto.Category

	err := validation.ValidateDto(c, &reqDto)
	if err != nil {
		return responses.JSONError(c, err.Error(), err)
	}

	resp, err := h.svc.CreateCategory(c.Request().Context(), &reqDto)
	if err != nil {
		return responses.JSONError(c, "failed to create category", err)
	}

	return responses.JSONSuccess(c, "created new category", resp)
}

// HandleGetCategories handles requests to fetch categories list.
func (h *ListingsHandler) HandleGetCategories(c echo.Context) error {
	resp, err := h.svc.GetCategories(c.Request().Context())
	if err != nil {
		return responses.JSONError(
			c,
			"failed to fetch categories",
			err,
			http.StatusInternalServerError,
		)
	}

	return responses.JSONSuccess(c, "fetched categories", resp)
}

// HandleCreateListing handles requests to create new listing.
func (h *ListingsHandler) HandleCreateListing(c echo.Context) error {
	return responses.JSONSuccess(c, "created new item", nil)
}

// HandleGetListings handles requests to get a list of listings.
func (h *ListingsHandler) HandleGetListings(c echo.Context) error {
	return responses.JSONSuccess(c, "fetched items", nil)
}
