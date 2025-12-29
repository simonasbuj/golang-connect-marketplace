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
	CreateUser(ctx context.Context, reqDto *dto.RegisterRequest) (*dto.User, error)
	GetUserByEmail(ctx context.Context, email string) (*dto.User, error)
	GetUserByID(ctx context.Context, id string) (*dto.User, error)
	GetUserByOAuthID(ctx context.Context, id string, provider dto.OAuthProvider) (*dto.User, error)
	SaveRefreshToken(ctx context.Context, token *dto.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*dto.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	// CreateOAuthUser(ctx context.Context, id string, provider dto.OAuthProvider) (*dto.User, error)
}

type repo struct {
	db *sqlx.DB
}

// New create new instance of auth/users repository.
func New(db *sqlx.DB) *repo { //nolint:revive
	return &repo{db: db}
}

func (r *repo) CreateUser(ctx context.Context, reqDto *dto.RegisterRequest) (*dto.User, error) {
	reqDto.ID = generate.ID("user")

	query := `
		INSERT INTO auth.users (id, email, password_hash, name, lastname, username) 
		VALUES (:id, :email, :password, :name, :lastname, :username)
		RETURNING id, email, name, lastname, username, role
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

func (r *repo) GetUserByEmail(ctx context.Context, email string) (*dto.User, error) {
	query := `
		SELECT id, email, password_hash, name, lastname, username, role 
		FROM auth.users 
		WHERE email = $1
	`

	var user dto.User

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("getting user by email from database: %w", err)
	}

	return &user, nil
}

func (r *repo) GetUserByID(ctx context.Context, id string) (*dto.User, error) {
	query := `
		SELECT id, email, password_hash, name, lastname, username, role 
		FROM auth.users 
		WHERE id = $1
	`

	var user dto.User

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("getting user by id from database: %w", err)
	}

	return &user, nil
}

func (r *repo) GetUserByOAuthID(ctx context.Context, id string, provider dto.OAuthProvider) (*dto.User, error) {
	query := `
		SELECT a.id, a.email, a.password_hash, a.name, a.lastname, a.username, a.role 
		FROM auth.users a
			LEFT JOIN auth.oauth_users o ON o.user_id = a.id
		WHERE o.provider_user_id = $1 and o.provider = $2
	`

	var user dto.User

	err := r.db.GetContext(ctx, &user, query, id, provider)
	if err != nil {
		return nil, fmt.Errorf("getting user by oatuh id from database: %w", err)
	}

	return &user, nil
}

func (r *repo) SaveRefreshToken(ctx context.Context, token *dto.RefreshToken) error {
	query := `INSERT INTO auth.refresh_tokens (token, user_id, expires_at) VALUES (:token,:user_id,:expires_at)`

	_, err := r.db.NamedExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("saving refresh token in database: %w", err)
	}

	return nil
}

func (r *repo) GetRefreshToken(ctx context.Context, tokenStr string) (*dto.RefreshToken, error) {
	query := `
		SELECT token, user_id, expires_at 
		FROM auth.refresh_tokens 
		WHERE token = $1
	`

	var token dto.RefreshToken

	err := r.db.GetContext(ctx, &token, query, tokenStr)
	if err != nil {
		return nil, fmt.Errorf("getting refresh token from database: %w", err)
	}

	return &token, nil
}

func (r *repo) DeleteRefreshToken(ctx context.Context, tokenStr string) error {
	q := `DELETE FROM auth.refresh_tokens WHERE token = $1`

	_, err := r.db.ExecContext(ctx, q, tokenStr)
	if err != nil {
		return fmt.Errorf("deleting refresh token from database: %w", err)
	}

	return nil
}
