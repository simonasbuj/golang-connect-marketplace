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
	AddListingImage(ctx context.Context, listingID, path string) (*dto.ListingImage, error)
	DeleteListingImage(ctx context.Context, req *dto.DeleteImageRequest) error
	UpdateListing(ctx context.Context, req *dto.UpdateListingRequest) (*dto.Listing, error)
	GetListings(ctx context.Context, req *dto.GetListingsRequest) ([]dto.Listing, error)
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
			c.title as "category_title",
			a.id as "seller.id",
			a.username as "seller.username",
			a.created_at as "seller.created_at",
			sa.id as "seller.seller_id",
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
			LEFT JOIN auth.users a on a.id = l.user_id
			LEFT JOIN payments.seller_accounts sa on sa.user_id = a.id
			LEFT JOIN listings.categories c on c.id = l.category_id
		WHERE l.id = $1
		GROUP BY l.id, a.id, sa.id, c.title
	`

	var listing dto.Listing

	err := r.db.GetContext(ctx, &listing, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("fetching listing by id from database: %w", err)
	}

	return &listing, nil
}

func (r *listingsRepo) AddListingImage(
	ctx context.Context,
	listingID, path string,
) (*dto.ListingImage, error) {
	id := generate.ID("img")

	query := `
		INSERT INTO listings.listings_images (id, listing_id, path) 
		VALUES ($1, $2, $3)
		RETURNING id, listing_id, path
	`

	var img dto.ListingImage

	err := r.db.GetContext(ctx, &img, query, id, listingID, path)
	if err != nil {
		return nil, fmt.Errorf("inserting listing image: %w", err)
	}

	return &img, nil
}

func (r *listingsRepo) DeleteListingImage(ctx context.Context, req *dto.DeleteImageRequest) error {
	query := `
		DELETE FROM listings.listings_images
		WHERE listing_id = $1 and id = $2
	`

	_, err := r.db.ExecContext(ctx, query, req.ListingID, req.ImageID)
	if err != nil {
		return fmt.Errorf("deleting listing image from database: %w", err)
	}

	return nil
}

func (r *listingsRepo) UpdateListing(
	ctx context.Context,
	req *dto.UpdateListingRequest,
) (*dto.Listing, error) {
	query := `
		UPDATE listings.listings
		SET
			category_id = COALESCE(NULLIF(:category_id, ''), category_id),
			title  = COALESCE(NULLIF(:title, ''), title),
			description  = COALESCE(NULLIF(:description, ''), description),
			price_in_cents  = COALESCE(NULLIF(:price_in_cents, 0), price_in_cents),
			currency  = COALESCE(NULLIF(:currency, ''), currency),
			status = COALESCE(:status, status),
			updated_at = NOW()
		WHERE id = :id
	`

	row, err := r.db.NamedQueryContext(ctx, query, req)
	if err != nil {
		return nil, fmt.Errorf("updating listing in database: %w", err)
	}

	defer func() { _ = row.Close() }()

	updatedListing, err := r.GetListingByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("fetching updated listing: %w", err)
	}

	return updatedListing, err
}

func (r *listingsRepo) GetListings(
	ctx context.Context,
	req *dto.GetListingsRequest,
) ([]dto.Listing, error) {
	query := `
		SELECT
			l.*,
			c.title as "category_title",
			a.id as "seller.id",
			a.username as "seller.username",
			a.created_at as "seller.created_at",
			sa.id as "seller.seller_id",
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
			LEFT JOIN auth.users a on a.id = l.user_id
			LEFT JOIN payments.seller_accounts sa on sa.user_id = a.id
			LEFT JOIN listings.categories c on c.id = l.category_id
		WHERE 
			l.status = 'open'
			AND ($3::text IS NULL OR c.title ILIKE '%' || $3 || '%')
		GROUP BY l.id, a.id, sa.id, c.title
		ORDER BY l.created_at DESC
		LIMIT $1 OFFSET $2;
	`

	listings := []dto.Listing{}

	err := r.db.SelectContext(
		ctx,
		&listings,
		query,
		req.Limit,
		req.Page*req.Limit,
		req.CategoryFilter,
	)
	if err != nil {
		return nil, fmt.Errorf("fetching listings from database: %w", err)
	}

	return listings, nil
}
