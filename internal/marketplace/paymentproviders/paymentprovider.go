package paymentproviders

import (
	"context"
	"golang-connect-marketplace/internal/marketplace/dto"
	"net/http"
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
	VerifySuccessWebhook(
		ctx context.Context,
		payload []byte,
		header http.Header,
	) (*dto.Payment, error)
	VerifyRefundWebhook(
		ctx context.Context,
		payload []byte,
		header http.Header,
	) (*dto.Payment, error)
}
