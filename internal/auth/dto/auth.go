// Package dto contains data transfer objects used for API requests/responses.
package dto

// RegisterRequest represents the payload sent when user is signing up.
type RegisterRequest struct {
	ID       string `json:"-"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=12"`
	Name     string `json:"name"     validate:"required"`
	Lastname string `json:"lastname" validate:"required"`
}

// User represents user object.
type User struct {
	ID       string `json:"id"       db:"id"`
	Email    string `json:"email"    db:"email"`
	Name     string `json:"name"     db:"name"`
	Lastname string `json:"lastname" db:"lastname"`
}
