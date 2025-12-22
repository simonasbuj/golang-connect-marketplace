// Package services implements the core business logic of listings and payments.
package services

import (
	"context"
	"fmt"
	authDto "golang-connect-marketplace/internal/auth/dto"
	"golang-connect-marketplace/internal/marketplace/dto"
	"golang-connect-marketplace/internal/marketplace/repos"
)

// ListingsService provides user and auth related operations.
type ListingsService struct {
	repo repos.ListingsRepo
}

// NewListingsService creates a new Service instance.
func NewListingsService(repo repos.ListingsRepo) *ListingsService {
	return &ListingsService{
		repo: repo,
	}
}

// CreateCategory handles logic for creating new category.
func (s *ListingsService) CreateCategory(
	ctx context.Context,
	req *dto.Category,
) (*dto.Category, error) {
	resp, err := s.repo.CreateCategory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("creating category: %w", err)
	}

	return resp, nil
}

// GetCategories handles logic for fetching categories list.
func (s *ListingsService) GetCategories(ctx context.Context) ([]dto.Category, error) {
	resp, err := s.repo.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching categories: %w", err)
	}

	return resp, nil
}

// CreateListing handles logic for creating a new listing.
func (s *ListingsService) CreateListing(
	ctx context.Context,
	userClaims *authDto.UserClaims,
	req *dto.Listing,
) (*dto.Listing, error) {
	req.UserID = userClaims.ID

	resp, err := s.repo.CreateListing(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("creating new listing: %w", err)
	}

	return resp, nil
}
