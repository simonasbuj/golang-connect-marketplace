// Package dto contains data transfer objects used for API requests/responses.
package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"time"
)

// Category represents category model.
type Category struct {
	ID          string `json:"id"          db:"id"`
	Title       string `json:"title"       db:"title"       validate:"required,max=30"`
	Description string `json:"description" db:"description"`
	Color       string `json:"color"       db:"color"       validate:"required,len=7"`
}

// ListingStatus represents the current lifecycle state of a marketplace listing.
type ListingStatus string

const (
	// ListingStatusOpen indicates that the listing is active and available.
	ListingStatusOpen ListingStatus = "open"
	// ListingStatusCanceled indicates that the listing was canceled by the seller.
	ListingStatusCanceled ListingStatus = "canceled"
	// ListingStatusSold indicates that the listing has been completed (sold).
	ListingStatusSold ListingStatus = "sold"
	// ListingStatusRefunded indicates that the listing has been refunded.
	ListingStatusRefunded ListingStatus = "refunded"
)

// Listing represents a marketplace listing created by a user.
type Listing struct {
	ID            string        `json:"id"             db:"id"`
	UserID        string        `json:"user_id"        db:"user_id"`
	CategoryID    string        `json:"category_id"    db:"category_id"    validate:"required"`
	CategoryTitle string        `json:"category_title" db:"category_title"`
	Title         string        `json:"title"          db:"title"          validate:"required,min=8,max=100"`
	Description   string        `json:"description"    db:"description"    validate:"required"`
	PriceInCents  int           `json:"price_in_cents" db:"price_in_cents" validate:"required,min=1000"`
	Currency      string        `json:"currency"       db:"currency"       validate:"required,len=3"`
	Seller        SellerAccount `json:"seller"         db:"seller"`
	Status        ListingStatus `json:"status"         db:"status"`
	Images        ListingImages `json:"images"         db:"images"`
	CreatedAt     time.Time     `json:"created_at"     db:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"     db:"updated_at"`
	DeletedAt     *time.Time    `json:"deleted_at"     db:"deleted_at"`
}

// AddImagesRequest represents payload sent when adding new images for a listing.
type AddImagesRequest struct {
	UserID      string                 `form:"-"      validate:"required"`
	ListingID   string                 `form:"-"      validate:"required,min=30"`
	FileHeaders []multipart.FileHeader `form:"images" validate:"required"`
}

// DeleteImageRequest represents payload sent when deleting an images from a listing.
type DeleteImageRequest struct {
	UserID    string `json:"-"        validate:"required"`
	ListingID string `json:"-"        validate:"required"`
	ImageID   string `json:"image_id" validate:"required,min=10"`
}

// ListingImage represents a single image belonging to a listing.
type ListingImage struct {
	ID        string `json:"id"         db:"id"`
	ListingID string `json:"listing_id" db:"listing_id"`
	Path      string `json:"path"       db:"path"`
}

// ListingImages represents a collection of listing images.
type ListingImages []ListingImage

// ErrInvalidListingImagesScanType is returned if scanning json into ListingImage fails.
var ErrInvalidListingImagesScanType = errors.New("invalid type for ListingImages scan")

// Scan implements sql.Scanner to decode JSON-aggregated images from the database.
func (li *ListingImages) Scan(value any) error {
	if value == nil {
		*li = ListingImages{}

		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("%w: %T", ErrInvalidListingImagesScanType, value)
	}

	err := json.Unmarshal(bytes, li)
	if err != nil {
		return fmt.Errorf("unmarshaling ListingImage dto: %w", err)
	}

	return nil
}

// UpdateListingRequest represents payload sent when adding updating a listing.
type UpdateListingRequest struct {
	ID           string         `json:"id"             db:"id"`
	CategoryID   string         `json:"category_id"    db:"category_id"`
	Title        string         `json:"title"          db:"title"          validate:"omitempty,min=8,max=100"`
	Description  string         `json:"description"    db:"description"`
	PriceInCents int            `json:"price_in_cents" db:"price_in_cents" validate:"omitempty,min=1000"`
	Currency     string         `json:"currency"       db:"currency"       validate:"omitempty,len=3"`
	Status       *ListingStatus `json:"status"         db:"status"`
}

// GetListingsRequest represents payload sent when fetching a list of listings.
type GetListingsRequest struct {
	Limit          int     `json:"limit"           validate:"omitempty,min=1,max=100" query:"limit"`
	Page           int     `json:"page"            validate:"omitempty,min=1"         query:"page"`
	CategoryFilter *string `json:"category_filter" validate:"-"                       query:"category_filter"`
	ListingFilter  *string `json:"listing_filter"  validate:"-"                       query:"listing_filter"`
}

// PaginationMeta represents pagination metadata sent back to the client.
type PaginationMeta struct {
	Limit          int     `json:"limit"`
	Page           int     `json:"page"`
	Total          int     `json:"total"`
	CategoryFilter *string `json:"category_filter"`
	ListingFilter  *string `json:"listing_filter"`
}

// GetListingsResponse represents payload sent back when fetching a list of listings.
type GetListingsResponse struct {
	Meta     PaginationMeta `json:"meta"`
	Listings []Listing      `json:"listings"`
}
