package dto

import (
	"time"
)

// Provider represents a payment provider enum.
type Provider string

const (
	// ProviderStripe represents the Stripe payment provider.
	ProviderStripe Provider = "stripe"
	// ProviderKlix represents the klix payment provider.
	ProviderKlix Provider = "klix"
	// ProviderPolar represents the polar.sh payment provider.
	ProviderPolar Provider = "polar.sh"
)

// SellerAcountLinkingSessionRequest represents payload sent when linking seller account.
type SellerAcountLinkingSessionRequest struct {
	UserID     string `json:"-"           validate:"required"`
	RefreshURL string `json:"refresh_url" validate:"required"`
	ReturnURL  string `json:"return_url"  validate:"required"`
}

// SellerAcountLinkingSessionResponse represents payload sent back when linking seller account.
type SellerAcountLinkingSessionResponse struct {
	SellerID string   `json:"seller_id"`
	URL      string   `json:"url"`
	Provider Provider `json:"provider"`
}

// SellerAccount represents a seller account.
type SellerAccount struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	Lastname  string    `db:"lastname"`
	Username  string    `db:"username"`
	SellerID  *string   `db:"seller_id"`
	CreatedAt time.Time `db:"created_at"`
}

// CheckoutSessionRequest represents payload sent when creating checkout session.
type CheckoutSessionRequest struct {
	BuyerID    string `json:"-"           validate:"required"`
	ListingID  string `json:"-"           validate:"required"`
	SuccessURL string `json:"success_url" validate:"required"`
	CancelURL  string `json:"cancel_url"  validate:"required"`
}

// CheckoutSessionResponse represents payload sent back when creating checkout session.
type CheckoutSessionResponse struct {
	URL string `json:"url"`
}
