// Package dto contains data transfer objects used for API requests/responses.
package dto

// Category represents category model.
type Category struct {
	ID          string `json:"id"          db:"id"`
	Title       string `json:"title"       db:"title"       validate:"required,max=30"`
	Description string `json:"description" db:"description"`
}
