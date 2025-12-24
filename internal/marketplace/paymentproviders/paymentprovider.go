package paymentproviders

import (
	"context"
	"golang-connect-marketplace/internal/marketplace/dto"
)

// PaymentProvider is an interface for payment-related operations.
type PaymentProvider interface {
	CreateAcountLinkingSession(
		ctx context.Context,
		req *dto.SellerAcountLinkingSessionRequest,
		user *dto.SellerAccount,
	) (*dto.SellerAcountLinkingSessionResponse, error)
	CreateAccountUpdateSession(
		ctx context.Context,
		req *dto.SellerAcountLinkingSessionRequest,
		user *dto.SellerAccount,
	) (*dto.SellerAcountLinkingSessionResponse, error)
}
