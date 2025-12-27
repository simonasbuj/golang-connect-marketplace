// Package services implements the core business logic of listings and payments.
package services

import (
	"context"
	"errors"
	"fmt"
	"golang-connect-marketplace/config"
	authDto "golang-connect-marketplace/internal/auth/dto"
	"golang-connect-marketplace/internal/marketplace/dto"
	"golang-connect-marketplace/internal/marketplace/repos"
	"golang-connect-marketplace/internal/marketplace/storage"
)

var (
	// ErrForbidden is returned when user is not allowed to do some actions.
	ErrForbidden = errors.New("user is forbidden to do this action")
	// ErrTooManyImages is returned when listing already has too many images.
	ErrTooManyImages = errors.New("listing has too many images")
	// ErrImageDoesntExist is returned when provided image id doesnt exist in listing.
	ErrImageDoesntExist = errors.New("provided image doesn't exist in this listing")
	// ErrListingIsNotOpen is returned when trying to update or pay for a listing that is not open anymore.
	ErrListingIsNotOpen = errors.New("listing is not open")
	// ErrUserIsNotSeller is returned user doesn't have seller account linked.
	ErrUserIsNotSeller = errors.New("user doesn't have seller account linked")
)

// ListingsService provides listing related operations bussines logic.
type ListingsService struct {
	repo    repos.ListingsRepo
	storage storage.Storage
	cfg     *config.StorageConfig
}

// NewListingsService creates a new Service instance.
func NewListingsService(
	repo repos.ListingsRepo,
	storage storage.Storage,
	cfg *config.StorageConfig,
) *ListingsService {
	return &ListingsService{
		repo:    repo,
		storage: storage,
		cfg:     cfg,
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

// AddImages handles logic for adding images to listing.
func (s *ListingsService) AddImages(
	ctx context.Context,
	req *dto.AddImagesRequest,
) (*dto.Listing, error) {
	listing, err := s.repo.GetListingByID(ctx, req.ListingID)
	if err != nil {
		return nil, fmt.Errorf("fetching listing: %w", err)
	}

	if listing.UserID != req.UserID {
		return nil, ErrForbidden
	}

	if len(listing.Images)+len(req.FileHeaders) > s.cfg.MaxImagesPerListing {
		return nil, ErrTooManyImages
	}

	for _, fh := range req.FileHeaders {
		path, err := s.storage.StoreImage(ctx, &fh, req.ListingID)
		if err != nil {
			return nil, fmt.Errorf("storing image: %w", err)
		}

		listingImage, err := s.repo.AddListingImage(ctx, req.ListingID, path)
		if err != nil {
			return nil, fmt.Errorf("inserting image: %w", err)
		}

		listing.Images = append(listing.Images, *listingImage)
	}

	return listing, nil
}

// DeleteImage handles logic for adding images to listing.
func (s *ListingsService) DeleteImage(
	ctx context.Context,
	req *dto.DeleteImageRequest,
) (*dto.Listing, error) {
	listing, err := s.repo.GetListingByID(ctx, req.ListingID)
	if err != nil {
		return nil, fmt.Errorf("fetching listing: %w", err)
	}

	if listing.Status != dto.ListingStatusOpen {
		return nil, ErrListingIsNotOpen
	}

	if listing.UserID != req.UserID {
		return nil, ErrForbidden
	}

	path := ""

	for i, img := range listing.Images {
		if img.ID == req.ImageID {
			listing.Images = append(listing.Images[:i], listing.Images[i+1:]...)
			path = img.Path

			break
		}
	}

	if path == "" {
		return nil, ErrImageDoesntExist
	}

	err = s.storage.DeleteImage(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("deleting image from storage: %w", err)
	}

	err = s.repo.DeleteListingImage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("deleting image from database: %w", err)
	}

	return listing, nil
}

// GetListingByID handles logic for fetching listing by id.
func (s *ListingsService) GetListingByID(
	ctx context.Context,
	listingID string,
) (*dto.Listing, error) {
	listing, err := s.repo.GetListingByID(ctx, listingID)
	if err != nil {
		return nil, fmt.Errorf("fetching listing: %w", err)
	}

	return listing, nil
}

// UpdateListing handles logic for updating a listing.
func (s *ListingsService) UpdateListing(
	ctx context.Context,
	req *dto.UpdateListingRequest,
	user *authDto.UserClaims,
) (*dto.Listing, error) {
	listing, err := s.repo.GetListingByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("fetching listing: %w", err)
	}

	if listing.UserID != user.ID && user.Role != authDto.UserRoleAdmin {
		return nil, ErrForbidden
	}

	if listing.Status != dto.ListingStatusOpen && user.Role != authDto.UserRoleAdmin {
		return nil, ErrListingIsNotOpen
	}

	updatedListing, err := s.repo.UpdateListing(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("updating listing: %w", err)
	}

	return updatedListing, nil
}

// GetListings handles logic for fetching a list of listing.
func (s *ListingsService) GetListings(
	ctx context.Context,
	req *dto.GetListingsRequest,
) (*dto.GetListingsResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}

	if req.Page < 0 {
		req.Page = 0
	}

	listings, err := s.repo.GetListings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("fetching listings: %w", err)
	}

	resp := &dto.GetListingsResponse{
		Meta: dto.PaginationMeta{
			Limit:          req.Limit,
			Page:           req.Page,
			CategoryFilter: req.CategoryFilter,
			ListingFilter:  req.ListingFilter,
			Total:          len(listings),
		},
		Listings: listings,
	}

	return resp, nil
}
