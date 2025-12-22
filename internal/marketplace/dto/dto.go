// Package dto contains data transfer objects used for API requests/responses.
package dto

import (
	"mime/multipart"
	"time"
)

// Category represents category model.
type Category struct {
	ID          string `json:"id"          db:"id"`
	Title       string `json:"title"       db:"title"       validate:"required,max=30"`
	Description string `json:"description" db:"description"`
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
	ID           string        `json:"id"             db:"id"`
	UserID       string        `json:"user_id"        db:"user_id"`
	CategoryID   string        `json:"category_id"    db:"category_id"    validate:"required"`
	Title        string        `json:"title"          db:"title"          validate:"required,min=8,max=100"`
	Description  string        `json:"description"    db:"description"    validate:"required"`
	PriceInCents int           `json:"price_in_cents" db:"price_in_cents" validate:"required,min=1"`
	Currency     string        `json:"currency"       db:"currency"       validate:"required,len=3"`
	Status       ListingStatus `json:"status"         db:"status"`
	CreatedAt    time.Time     `json:"created_at"     db:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"     db:"updated_at"`
}

// AddImagesRequest represents payload sent when adding new images for a listing.
type AddImagesRequest struct {
	UserID      string                 `validate:"required"`
	ListingID   string                 `validate:"required,min=30"`
	FileHeaders []multipart.FileHeader `validate:"required"        form:"images"`
}
