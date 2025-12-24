package paymentproviders

import (
	"context"
	"golang-connect-marketplace/internal/marketplace/dto"
)

type PaymentProvider interface {
	CreateAcountLinkingSession(
		ctx context.Context,
		req *dto.SellerAcountLinkingSessionRequest,
	) (*dto.SellerAcountLinkingSessionResponse, error)
}
