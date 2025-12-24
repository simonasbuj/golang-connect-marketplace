package dto

import "time"

type SellerAcountLinkingSessionRequest struct {
	UserID     string `json:"-"           validate:"required"`
	RefreshURL string `json:"refresh_url" validate:"required"`
	ReturnURL  string `json:"return_url"  validate:"required"`
}

type SellerAcountLinkingSessionResponse struct {
	SellerID string `json:"seller_id"`
	URL      string `json:"url"`
}

type SellerAccount struct {
	UserID    string
	Email     string
	Name      string
	Lastname  string
	Username  string
	SellerID  string
	CreatedAt time.Time
}
