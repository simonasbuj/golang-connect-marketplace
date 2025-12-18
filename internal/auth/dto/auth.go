// Package dto contains data transfer objects used for API requests/responses.
package dto

import "time"

// RegisterRequest represents the payload sent when user is signing up.
type RegisterRequest struct {
	ID       string `json:"-"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=12"`
	Name     string `json:"name"     validate:"required,max=100"`
	Lastname string `json:"lastname" validate:"required,max=100"`
	Username string `json:"username" validate:"required,max=40"`
}

// UserRole represents the role of a user in the database.
type UserRole string

const (
	// UserRoleAdmin user with full permissions.
	UserRoleAdmin UserRole = "admin"
	// UserRoleCustomer customer user.
	UserRoleCustomer UserRole = "customer"
)

// User represents user object.
type User struct {
	ID           string   `json:"id"       db:"id"`
	Email        string   `json:"email"    db:"email"`
	Name         string   `json:"name"     db:"name"`
	Lastname     string   `json:"lastname" db:"lastname"`
	Uername      string   `json:"username" db:"username"`
	Role         UserRole `json:"role"     db:"role"`
	PasswordHash string   `json:"-"        db:"password_hash"`
}

// LoginRequest represents the payload sent when user is loging in.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=12"`
}

// LoginResponse represents the payload sent back to user after loging in.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UserClaims represents user claim saved in jwt token.
type UserClaims struct {
	ID   string   `json:"id"   db:"id"`
	Role UserRole `json:"role" db:"role"`
}

// RefreshToken represents refresh token row in database.
type RefreshToken struct {
	Token     string    `db:"token"`
	UserID    string    `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
}
