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

const categoryIDLength = 26

// ListingsRepo defines methods for accessing and managing listings data.
type ListingsRepo interface {
	CreateCategory(ctx context.Context, req *dto.Category) (*dto.Category, error)
	GetCategories(ctx context.Context) ([]dto.Category, error)
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
	req.ID = generate.ID("cat", categoryIDLength)

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
	var categories []dto.Category

	query := `
		SELECT id, title, description
		FROM listings.categories
		WHERE deleted_at IS NULL
		ORDER BY title
	`

	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("fetching categories fropm database: %w", err)
	}

	return categories, nil
}
