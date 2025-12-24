package repos

import (
	"context"
	"golang-connect-marketplace/internal/marketplace/dto"
)

// PaymentsRepo defines methods for accessing and managing payments data.
type PaymentsRepo interface {
	GetSellerInfoByID(ctx context.Context, userID string) (*dto.SellerAccount, error)
}
