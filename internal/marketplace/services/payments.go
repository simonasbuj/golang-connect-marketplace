package services

import (
	"context"
	"fmt"
	"golang-connect-marketplace/internal/marketplace/dto"
	"golang-connect-marketplace/internal/marketplace/paymentproviders"
	"golang-connect-marketplace/internal/marketplace/repos"
)

// PaymentsService provides payments related operations bussines logic.
type PaymentsService struct {
	provider paymentproviders.PaymentProvider
	repo     repos.PaymentsRepo
}

// NewPaymentsService returns an instance of PaymentsService.
func NewPaymentsService(
	provider paymentproviders.PaymentProvider,
	repo repos.PaymentsRepo,
) *PaymentsService {
	return &PaymentsService{
		provider: provider,
		repo:     repo,
	}
}

// LinkSellerAccount handles bussines logic for linking users to seller accounts.
func (s *PaymentsService) LinkSellerAccount(
	ctx context.Context,
	req *dto.SellerAcountLinkingSessionRequest,
) (*dto.SellerAcountLinkingSessionResponse, error) {
	user, err := s.repo.GetSellerInfoByID(ctx, req.UserID)
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

	_, err = s.repo.UpdateSellerID(ctx, req.UserID, resp.SellerID)
	if err != nil {
		return nil, fmt.Errorf("updating seller id: %w", err)
	}

	return resp, nil
}
