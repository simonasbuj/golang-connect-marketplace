package services

import (
	"context"
	"fmt"
	"golang-connect-marketplace/internal/marketplace/dto"
	"golang-connect-marketplace/internal/marketplace/paymentproviders"
	"golang-connect-marketplace/internal/marketplace/repos"
	"net/http"
)

const (
	minimumFee     = 100
	feePercent     = 4
	percentDivisor = 100
)

// PaymentsService provides payments related operations bussines logic.
type PaymentsService struct {
	provider     paymentproviders.PaymentProvider
	paymentsRepo repos.PaymentsRepo
	listingsRepo repos.ListingsRepo
}

// NewPaymentsService returns an instance of PaymentsService.
func NewPaymentsService(
	provider paymentproviders.PaymentProvider,
	paymentsRepo repos.PaymentsRepo,
	listingsRepo repos.ListingsRepo,
) *PaymentsService {
	return &PaymentsService{
		provider:     provider,
		paymentsRepo: paymentsRepo,
		listingsRepo: listingsRepo,
	}
}

// LinkSellerAccount handles bussines logic for linking users to seller accounts.
func (s *PaymentsService) LinkSellerAccount(
	ctx context.Context,
	req *dto.SellerAcountLinkingSessionRequest,
) (*dto.SellerAcountLinkingSessionResponse, error) {
	user, err := s.paymentsRepo.GetSellerInfoByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("fetching seller info: %w", err)
	}

	if user.SellerID != nil {
		resp, err := s.provider.CreateAccountUpdateSession(ctx, req, user)
		if err != nil {
			return nil, fmt.Errorf("creating seller account update session: %w", err)
		}

		return resp, nil
	}

	resp, err := s.provider.CreateAcountLinkingSession(ctx, req, user)
	if err != nil {
		return nil, fmt.Errorf("creating seller account linking session: %w", err)
	}

	_, err = s.paymentsRepo.UpdateSellerID(ctx, req.UserID, resp.SellerID, resp.Provider)
	if err != nil {
		return nil, fmt.Errorf("updating seller id: %w", err)
	}

	return resp, nil
}

// CreateCheckoutSession handles bussines logic for creating checkout session for a listing.
func (s *PaymentsService) CreateCheckoutSession(
	ctx context.Context,
	req *dto.CheckoutSessionRequest,
) (*dto.CheckoutSessionResponse, error) {
	listing, err := s.listingsRepo.GetListingByID(ctx, req.ListingID)
	if err != nil {
		return nil, fmt.Errorf("fetching listing while creating checkout session: %w", err)
	}

	if listing.Status != dto.ListingStatusOpen {
		return nil, ErrListingIsNotOpen
	}

	if listing.Seller.SellerID == nil {
		return nil, ErrUserIsNotSeller
	}

	fee := s.calculateFee(listing)

	resp, err := s.provider.CreateCheckoutSession(ctx, req, listing, fee)
	if err != nil {
		return nil, fmt.Errorf("creating checkout session: %w", err)
	}

	return resp, nil
}

// HandleSuccessWebhook handles bussines logic for payment success webhook.
func (s *PaymentsService) HandleSuccessWebhook(
	ctx context.Context,
	payload []byte,
	header http.Header,
) (*dto.Payment, error) {
	payment, err := s.provider.VerifySuccessWebhook(ctx, payload, header)
	if err != nil {
		return nil, fmt.Errorf("verifying payment success webhook: %w", err)
	}

	_, err = s.paymentsRepo.SavePayment(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("saving payment: %w", err)
	}

	return payment, nil
}

// HandleRefundWebhook handles bussines logic for payment refund webhook.
func (s *PaymentsService) HandleRefundWebhook(
	ctx context.Context,
	payload []byte,
	header http.Header,
) (*dto.Payment, error) {
	payment, err := s.provider.VerifyRefundWebhook(ctx, payload, header)
	if err != nil {
		return nil, fmt.Errorf("verifying payment refunded webhook: %w", err)
	}

	ref, err := s.paymentsRepo.RefundPayment(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("refunding payments: %w", err)
	}

	return ref, nil
}

func (s *PaymentsService) calculateFee(listing *dto.Listing) int64 {
	fee := listing.PriceInCents*feePercent/percentDivisor + minimumFee

	return int64(fee)
}
