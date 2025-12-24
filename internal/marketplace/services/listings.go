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

	if listing.UserID != req.UserID {
		return nil, ErrForbidden
	}

	doesImageExist := false
	path := ""

	for i, img := range listing.Images {
		if img.ID == req.ImageID {
			listing.Images = append(listing.Images[:i], listing.Images[i+1:]...)
			doesImageExist = true
			path = img.Path

			break
		}
	}

	if !doesImageExist {
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
