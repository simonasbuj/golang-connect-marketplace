package services

import (
	"context"
	"fmt"
	"golang-connect-marketplace/internal/marketplace/dto"
	"golang-connect-marketplace/internal/marketplace/paymentproviders"
)

// PaymentsService provides payments related operations bussines logic.
type PaymentsService struct {
	provider paymentproviders.PaymentProvider
}

// NewPaymentsService returns an instance of PaymentsService.
func NewPaymentsService(provider paymentproviders.PaymentProvider) *PaymentsService {
	return &PaymentsService{
		provider: provider,
	}
}

// LinkSellerAccount handles bussines logic for linking users to seller accounts.
func (s *PaymentsService) LinkSellerAccount(
	ctx context.Context,
	req *dto.SellerAcountLinkingSessionRequest,
) (*dto.SellerAcountLinkingSessionResponse, error) {
	resp, err := s.provider.CreateAcountLinkingSession(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("creating seller account linking session: %w", err)
	}

	return resp, nil
}
