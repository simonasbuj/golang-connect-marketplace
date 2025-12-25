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
	CreateCheckoutSession(
		ctx context.Context,
		req *dto.CheckoutSessionRequest,
		listing *dto.Listing,
		feeAmount int64,
	) (*dto.CheckoutSessionResponse, error)
}
