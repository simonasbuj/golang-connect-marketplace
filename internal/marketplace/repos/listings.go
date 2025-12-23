// Package repos provides data access implementations for interacting with databases
package repos

import (
	"context"
	"errors"
	"fmt"
	"golang-connect-marketplace/internal/marketplace/dto"
	"golang-connect-marketplace/pkg/generate"

	"github.com/jmoiron/sqlx"
)

// ErrNoRowsReturned is returned if database query returns no results.
var ErrNoRowsReturned = errors.New("no rows returned")

// ListingsRepo defines methods for accessing and managing listings data.
type ListingsRepo interface {
	CreateCategory(ctx context.Context, req *dto.Category) (*dto.Category, error)
	GetCategories(ctx context.Context) ([]dto.Category, error)
	CreateListing(ctx context.Context, req *dto.Listing) (*dto.Listing, error)
	CheckIfUserOwnsListing(ctx context.Context, listingID, userID string) error
	GetListingByID(ctx context.Context, listingID string) (*dto.Listing, error)
}

type listingsRepo struct {
	db *sqlx.DB
}

// NewListingsRepo create new instance of lisitngs repository.
func NewListingsRepo(db *sqlx.DB) *listingsRepo { //nolint:revive
	return &listingsRepo{db: db}
}

func (r *listingsRepo) CreateCategory(
	ctx context.Context,
	req *dto.Category,
) (*dto.Category, error) {
	req.ID = generate.ID("cat")

	query := `
		INSERT INTO listings.categories (id, title, description) 
		VALUES (:id, :title, :description)
		RETURNING id, title, description
	`

	row, err := r.db.NamedQueryContext(ctx, query, req)
	if err != nil {
		return nil, fmt.Errorf("inserting category into database: %w", err)
	}

	defer func() { _ = row.Close() }()

	ok := row.Next()
	if !ok {
		return nil, ErrNoRowsReturned
	}

	var created dto.Category

	err = row.StructScan(&created)
	if err != nil {
		return nil, fmt.Errorf("scanning inserted category row into struct: %w", err)
	}

	return &created, nil
}

func (r *listingsRepo) GetCategories(ctx context.Context) ([]dto.Category, error) {
	categories := []dto.Category{}

	query := `
		SELECT id, title, description
		FROM listings.categories
		WHERE deleted_at IS NULL
		ORDER BY title
	`

	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("fetching categories from database: %w", err)
	}

	return categories, nil
}

func (r *listingsRepo) CreateListing(ctx context.Context, req *dto.Listing) (*dto.Listing, error) {
	req.ID = generate.ID("item")

	query := `
		INSERT INTO listings.listings (id, user_id, category_id, title, description, price_in_cents, currency) 
		VALUES (:id, :user_id, :category_id, :title, :description, :price_in_cents, :currency)
		RETURNING id, user_id, category_id, title, description, price_in_cents, currency, status, created_at, updated_at
	`

	row, err := r.db.NamedQueryContext(ctx, query, req)
	if err != nil {
		return nil, fmt.Errorf("inserting new listing to database: %w", err)
	}

	defer func() { _ = row.Close() }()

	ok := row.Next()
	if !ok {
		return nil, ErrNoRowsReturned
	}

	var resp dto.Listing

	err = row.StructScan(&resp)
	if err != nil {
		return nil, fmt.Errorf("scanning inserted listing row into struct: %w", err)
	}

	return &resp, nil
}

func (r *listingsRepo) CheckIfUserOwnsListing(ctx context.Context, listingID, userID string) error {
	query := `
		SELECT 1
		FROM listings.listings
		WHERE id = $1 and user_id = $2
	`

	var exists int

	err := r.db.GetContext(ctx, &exists, query, listingID, userID)
	if err != nil {
		return fmt.Errorf("confirming if user owns listing in database: %w", err)
	}

	return nil
}

func (r *listingsRepo) GetListingByID(ctx context.Context, listingID string) (*dto.Listing, error) {
	query := `
		SELECT
			l.*,
			COALESCE(
				json_agg(
					json_build_object(
						'id', i.id,
						'listing_id', i.listing_id,
						'path', i.path
					)
				) FILTER (WHERE i.id IS NOT NULL),
				'[]'
			) AS images
		FROM listings.listings l
		LEFT JOIN listings.listings_images i ON i.listing_id = l.id
		WHERE l.id = $1
		GROUP BY l.id
	`

	var listing dto.Listing

	err := r.db.GetContext(ctx, &listing, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("fetching listing by id from database: %w", err)
	}

	return &listing, nil
}
