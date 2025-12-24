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
	UpdateSellerID(ctx context.Context, userID, sellerID string) (*dto.SellerAccount, error)
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
		SELECT id, email, name, lastname, username, seller_id, created_at
		FROM auth.users
		WHERE id = $1
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
) (*dto.SellerAccount, error) {
	query := `
		UPDATE auth.users
		SET seller_id = $2
		WHERE id = $1
		RETURNING id, email, name, lastname, username, seller_id, created_at
	`

	var resp dto.SellerAccount

	err := r.db.GetContext(ctx, &resp, query, userID, sellerID)
	if err != nil {
		return nil, fmt.Errorf("updating seller_id for user in database: %w", err)
	}

	return &resp, nil
}
