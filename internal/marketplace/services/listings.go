// Package services implements the core business logic of listings and payments.
package services

// ListingsService provides user and auth related operations.
type ListingsService struct{}

// NewListingsService creates a new Service instance.
func NewListingsService() *ListingsService {
	return &ListingsService{}
}
