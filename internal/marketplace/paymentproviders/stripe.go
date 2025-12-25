// Package paymentproviders provides implementations of the PaymentProvider
package paymentproviders

import (
	"context"
	"fmt"
	"golang-connect-marketplace/internal/marketplace/dto"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/account"
	"github.com/stripe/stripe-go/v84/accountlink"
	"github.com/stripe/stripe-go/v84/checkout/session"
)

type stripePaymentProvider struct {
	webhookSecret string
}

const (
	metadataKeyListingID = "listing_id"
	metadataKeyBuyerID   = "buyer_id"
)

// NewStripePaymentProvider returns stripePaymentProvider which implements the PaymentProvider interface using stripe.
func NewStripePaymentProvider(
	secretKey, webhookSecret string,
) *stripePaymentProvider { //nolint:revive
	if secretKey == "" || webhookSecret == "" {
		panic("secretKey and webhookSecret are required for stripePaymentProvider")
	}

	stripe.Key = secretKey

	return &stripePaymentProvider{
		webhookSecret: webhookSecret,
	}
}

func (p *stripePaymentProvider) CreateAcountLinkingSession(
	_ context.Context,
	req *dto.SellerAcountLinkingSessionRequest,
	user *dto.SellerAccount,
) (*dto.SellerAcountLinkingSessionResponse, error) {
	params := &stripe.AccountParams{
		Type:  stripe.String(stripe.AccountTypeExpress),
		Email: stripe.String(user.Email),
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
		Provider: dto.ProviderStripe,
	}

	return resp, nil
}

func (p *stripePaymentProvider) CreateAccountUpdateSession(
	_ context.Context,
	req *dto.SellerAcountLinkingSessionRequest,
	user *dto.SellerAccount,
) (*dto.SellerAcountLinkingSessionResponse, error) {
	linkParams := &stripe.AccountLinkParams{
		Account:    stripe.String(*user.SellerID),
		RefreshURL: &req.RefreshURL,
		ReturnURL:  &req.ReturnURL,
		Type:       stripe.String("account_onboarding"),
	}

	link, err := accountlink.New(linkParams)
	if err != nil {
		return nil, fmt.Errorf("creating stripe account update session url: %w", err)
	}

	resp := &dto.SellerAcountLinkingSessionResponse{
		SellerID: *user.SellerID,
		URL:      link.URL,
		Provider: dto.ProviderStripe,
	}

	return resp, nil
}

func (p *stripePaymentProvider) CreateCheckoutSession(
	_ context.Context,
	req *dto.CheckoutSessionRequest,
	listing *dto.Listing,
	feeAmount int64,
) (*dto.CheckoutSessionResponse, error) {
	checkoutItemName := fmt.Sprintf("Buying %s from @%s", listing.Title, listing.Seller.Username)

	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(stripe.CheckoutSessionModePayment),
		SuccessURL: stripe.String(req.SuccessURL),
		CancelURL:  stripe.String(req.CancelURL),

		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			ApplicationFeeAmount: stripe.Int64(feeAmount),
			TransferData: &stripe.CheckoutSessionPaymentIntentDataTransferDataParams{
				Destination: stripe.String(*listing.Seller.SellerID),
			},
			Metadata: map[string]string{
				metadataKeyListingID: listing.ID,
				metadataKeyBuyerID:   req.BuyerID,
			},
		},

		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(listing.Currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(checkoutItemName),
					},
					UnitAmount: stripe.Int64(int64(listing.PriceInCents)),
				},
				Quantity: stripe.Int64(1),
			},
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(listing.Currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("marketplace fee"),
					},
					UnitAmount: stripe.Int64(feeAmount),
				},
				Quantity: stripe.Int64(1),
			},
		},
	}

	s, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("craeting stripe checkout session: %w", err)
	}

	resp := &dto.CheckoutSessionResponse{
		URL: s.URL,
	}

	return resp, nil
}
