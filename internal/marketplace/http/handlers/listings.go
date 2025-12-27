// Package handlers defines HTTP handlers for application marketplace endpoints.
package handlers

import (
	"errors"
	"golang-connect-marketplace/internal/auth/middleware"
	"golang-connect-marketplace/internal/marketplace/dto"
	"golang-connect-marketplace/internal/marketplace/services"
	r "golang-connect-marketplace/pkg/responses"
	"golang-connect-marketplace/pkg/validation"
	"net/http"

	"github.com/labstack/echo/v4"
)

const listingIDParamName = "listing_id"

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
		return r.JSONError(c, err.Error(), err)
	}

	resp, err := h.svc.CreateCategory(c.Request().Context(), &reqDto)
	if err != nil {
		return r.JSONError(c, "failed to create category", err)
	}

	return r.JSONSuccess(c, "created new category", resp)
}

// HandleGetCategories handles requests to fetch categories list.
func (h *ListingsHandler) HandleGetCategories(c echo.Context) error {
	resp, err := h.svc.GetCategories(c.Request().Context())
	if err != nil {
		return r.JSONError(c, "failed to fetch categories", err, http.StatusInternalServerError)
	}

	return r.JSONSuccess(c, "fetched categories", resp)
}

// HandleCreateListing handles requests to create new listing.
func (h *ListingsHandler) HandleCreateListing(c echo.Context) error {
	userClaims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	var reqDto dto.Listing

	err = validation.ValidateDto(c, &reqDto)
	if err != nil {
		return r.JSONError(c, err.Error(), err)
	}

	resp, err := h.svc.CreateListing(c.Request().Context(), userClaims, &reqDto)
	if err != nil {
		return r.JSONError(c, "failed to create listing", err, http.StatusInternalServerError)
	}

	return r.JSONSuccess(c, "created new listing", resp)
}

// HandleAddImages handles uploading images for a listing.
func (h *ListingsHandler) HandleAddImages(c echo.Context) error {
	userClaims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	listingID := c.Param(listingIDParamName)

	var reqDto dto.AddImagesRequest

	reqDto.UserID = userClaims.ID
	reqDto.ListingID = listingID

	err = validation.ValidateDto(c, &reqDto)
	if err != nil {
		return r.JSONError(c, err.Error(), err)
	}

	resp, err := h.svc.AddImages(c.Request().Context(), &reqDto)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			return r.JSONError(c, "forbidden", err, http.StatusForbidden)
		}

		if errors.Is(err, services.ErrTooManyImages) {
			return r.JSONError(c, "listing has too many images", err)
		}

		return r.JSONError(c, "failed to add images", err, http.StatusInternalServerError)
	}

	return r.JSONSuccess(c, "added images to listing", resp)
}

// HandleDeleteImages handles deleting image from a listing.
func (h *ListingsHandler) HandleDeleteImages(c echo.Context) error {
	userClaims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	listingID := c.Param(listingIDParamName)

	var reqDto dto.DeleteImageRequest

	reqDto.UserID = userClaims.ID
	reqDto.ListingID = listingID

	err = validation.ValidateDto(c, &reqDto)
	if err != nil {
		return r.JSONError(c, err.Error(), err)
	}

	resp, err := h.svc.DeleteImage(c.Request().Context(), &reqDto)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			return r.JSONError(c, "forbidden", err, http.StatusForbidden)
		}

		return r.JSONError(c, "failed to delete image", err)
	}

	return r.JSONSuccess(c, "deleted image from listing", resp)
}

// HandleGetListings handles requests to get a list of listings.
func (h *ListingsHandler) HandleGetListings(c echo.Context) error {
	var reqDto dto.GetListingsRequest

	err := validation.ValidateDto(c, &reqDto)
	if err != nil {
		return r.JSONError(c, err.Error(), err)
	}

	resp, err := h.svc.GetListings(c.Request().Context(), &reqDto)
	if err != nil {
		return r.JSONError(c, "failed to fetch listings", err)
	}

	return r.JSONSuccess(c, "fetched listings", resp)
}

// HandleGetListing handles requests to get a listing by id.
func (h *ListingsHandler) HandleGetListing(c echo.Context) error {
	listingID := c.Param(listingIDParamName)

	resp, err := h.svc.GetListingByID(c.Request().Context(), listingID)
	if err != nil {
		return r.JSONError(c, "failed to fetch listing", err, http.StatusInternalServerError)
	}

	return r.JSONSuccess(c, "fetched listing", resp)
}

// HandleUpdateListing handles requests to update a listing.
func (h *ListingsHandler) HandleUpdateListing(c echo.Context) error {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	listingID := c.Param(listingIDParamName)

	var reqDto dto.UpdateListingRequest

	reqDto.ID = listingID

	err = validation.ValidateDto(c, &reqDto)
	if err != nil {
		return r.JSONError(c, err.Error(), err)
	}

	resp, err := h.svc.UpdateListing(c.Request().Context(), &reqDto, user)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			return r.JSONError(c, "forbidden", err, http.StatusForbidden)
		}

		return r.JSONError(c, "failed to update listing", err)
	}

	return r.JSONSuccess(c, "updated listing", resp)
}
