package dto

import "time"

// SellerAcountLinkingSessionRequest represents payload sent when linking seller account.
type SellerAcountLinkingSessionRequest struct {
	UserID     string `json:"-"           validate:"required"`
	RefreshURL string `json:"refresh_url" validate:"required"`
	ReturnURL  string `json:"return_url"  validate:"required"`
}

// SellerAcountLinkingSessionResponse represents payload sent back when linking seller account.
type SellerAcountLinkingSessionResponse struct {
	SellerID string `json:"seller_id"`
	URL      string `json:"url"`
}

// SellerAccount represents a seller account.
type SellerAccount struct {
	UserID    string
	Email     string
	Name      string
	Lastname  string
	Username  string
	SellerID  string
	CreatedAt time.Time
}
