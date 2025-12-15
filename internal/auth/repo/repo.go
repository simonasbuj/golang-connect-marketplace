// Package repo provides data access implementations for interacting with databases
package repo

import (
	"context"
	"errors"
	"fmt"
	"golang-connect-marketplace/internal/auth/dto"
	"golang-connect-marketplace/pkg/generate"

	"github.com/jmoiron/sqlx"
)

// ErrNoRowsReturned is returned if database query returns no results.
var ErrNoRowsReturned = errors.New("no rows returned")

// Repo defines methods for accessing and managing user data.
type Repo interface {
	Create(ctx context.Context, reqDto *dto.RegisterRequest) (*dto.User, error)
}

type repo struct {
	db *sqlx.DB
}

// New create new instance of auth/users repository.
func New(db *sqlx.DB) *repo { //nolint:revive
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, reqDto *dto.RegisterRequest) (*dto.User, error) {
	reqDto.ID = generate.ID("user")

	query := `
		INSERT INTO auth.users (id, email, password_hash, name, lastname) 
		VALUES (:id, :email, :password, :name, :lastname)
		RETURNING id, email, name, lastname
	`

	row, err := r.db.NamedQueryContext(ctx, query, reqDto)
	if err != nil {
		return nil, fmt.Errorf("inserting user into database: %w", err)
	}

	defer func() { _ = row.Close() }()

	if !row.Next() {
		return nil, ErrNoRowsReturned
	}

	var user dto.User

	err = row.StructScan(&user)
	if err != nil {
		return nil, fmt.Errorf("scanning results into user dto: %w", err)
	}

	return &user, nil
}
