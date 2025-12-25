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
	ID        string    `json:"id,omitempty"         db:"id"`
	Email     string    `json:"email,omitempty"      db:"email"`
	Name      string    `json:"name,omitempty"       db:"name"`
	Lastname  string    `json:"lastname,omitempty"   db:"lastname"`
	Username  string    `json:"username,omitempty"   db:"username"`
	SellerID  *string   `json:"seller_id,omitempty"  db:"seller_id"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
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
