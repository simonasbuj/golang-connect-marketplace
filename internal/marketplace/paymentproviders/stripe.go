// package paymentproviders
package paymentproviders

import (
	"context"
	"fmt"
	"golang-connect-marketplace/internal/marketplace/dto"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/account"
	"github.com/stripe/stripe-go/v84/accountlink"
)

type stripePaymentProvider struct {
	webhookSecret string
}

func NewStripePaymentProvider(secretKey, webhookSecret string) *stripePaymentProvider {
	if secretKey == "" || webhookSecret == "" {
		panic("secretKey and webhookSecret are required for stripePaymentProvider")
	}

	stripe.Key = secretKey

	return &stripePaymentProvider{
		webhookSecret: webhookSecret,
	}
}

func (p *stripePaymentProvider) CreateAcountLinkingSession(
	ctx context.Context,
	req *dto.SellerAcountLinkingSessionRequest,
) (*dto.SellerAcountLinkingSessionResponse, error) {
	params := &stripe.AccountParams{
		Type: stripe.String(stripe.AccountTypeExpress),
	}

	acc, err := account.New(params)
	if err != nil {
		return nil, fmt.Errorf("creating new stripe seller account id: %w", err)
	}

	linkParams := &stripe.AccountLinkParams{
		Account:    stripe.String(acc.ID),
		RefreshURL: &req.RefreshURL,
		ReturnURL:  &req.ReturnURL,
		Type:       stripe.String("account_onboarding"),
	}

	link, err := accountlink.New(linkParams)
	if err != nil {
		return nil, fmt.Errorf("creating new stripe linking session url: %w", err)
	}

	resp := &dto.SellerAcountLinkingSessionResponse{
		SellerID: acc.ID,
		URL:      link.URL,
	}

	return resp, nil
}
