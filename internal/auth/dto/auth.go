// Package dto contains data transfer objects used for API requests/responses.
package dto

// RegisterRequest represents the payload sent when user is signing up.
type RegisterRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=12"`
	Name     string `json:"name"     validate:"required"`
	Lastname string `json:"lastname" validate:"required"`
}

// User represents user object.
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
}
