package repos

import (
	"context"
	"fmt"
	"golang-connect-marketplace/internal/marketplace/dto"

	"github.com/jmoiron/sqlx"
)

// PaymentsRepo defines methods for accessing and managing payments data.
type PaymentsRepo interface {
	GetSellerInfoByID(ctx context.Context, userID string) (*dto.SellerAccount, error)
	UpdateSellerID(
		ctx context.Context,
		userID, sellerID string,
		provider dto.Provider,
	) (*dto.SellerAccount, error)
	SavePayment(ctx context.Context, payment *dto.Payment) (*dto.Payment, error)
}

type paymentsRepo struct {
	db *sqlx.DB
}

// NewPaymentsRepo create new instance of payments repository.
func NewPaymentsRepo(db *sqlx.DB) *paymentsRepo { //nolint:revive
	return &paymentsRepo{db: db}
}

func (r *paymentsRepo) GetSellerInfoByID(
	ctx context.Context,
	userID string,
) (*dto.SellerAccount, error) {
	query := `
		SELECT a.id, a.email, a.name, a.lastname, a.username, s.id as seller_id, a.created_at
		FROM auth.users a
			LEFT JOIN payments.seller_accounts s ON a.id = s.user_id
		WHERE a.id = $1
	`

	var resp dto.SellerAccount

	err := r.db.GetContext(ctx, &resp, query, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching user from database: %w", err)
	}

	return &resp, nil
}

func (r *paymentsRepo) UpdateSellerID(
	ctx context.Context,
	userID, sellerID string,
	provider dto.Provider,
) (*dto.SellerAccount, error) {
	query := `
		WITH insert_stmt as (
			INSERT INTO payments.seller_accounts (id, user_id, provider)
			VALUES($2, $1, $3)
		)
		SELECT a.id, a.email, a.name, a.lastname, a.username, s.id as seller_id, a.created_at
		FROM auth.users a
			LEFT JOIN payments.seller_accounts s ON a.id = s.user_id
		WHERE a.id = $1
	`

	var resp dto.SellerAccount

	err := r.db.GetContext(ctx, &resp, query, userID, sellerID, provider)
	if err != nil {
		return nil, fmt.Errorf("updating seller_id for user in database: %w", err)
	}

	return &resp, nil
}

func (r *paymentsRepo) SavePayment(
	ctx context.Context,
	payment *dto.Payment,
) (*dto.Payment, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("starting db transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	insertPaymentQ := `
		INSERT INTO payments.payments 
			(id, listing_id, buyer_id, provider_payment_id, provider, amount_in_cents, fee_amount_in_cents, currency) 
		VALUES 
			(:id, :listing_id, :buyer_id, :provider_payment_id, :provider, :amount_in_cents, :fee_amount_in_cents, :currency) 
		RETURNING *
	`

	_, err = tx.NamedExecContext(ctx, insertPaymentQ, payment)
	if err != nil {
		return nil, fmt.Errorf("inserting payment into database: %w", err)
	}

	updateListingQ := `
		UPDATE listings.listings SET status = 'sold' WHERE id = $1 
	`

	_, err = tx.ExecContext(ctx, updateListingQ, payment.ListingID)
	if err != nil {
		return nil, fmt.Errorf("setting listing status to sold: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf(
			"committing transcation for saving payment + uptading listing: %w",
			err,
		)
	}

	return payment, nil
}
