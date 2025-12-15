// Package service implements the core business logic of authentication and user management.
package service

import (
	"context"
	"fmt"
	"golang-connect-marketplace/internal/auth/dto"
	"golang-connect-marketplace/internal/auth/repo"

	"golang.org/x/crypto/bcrypt"
)

// Service provides user and auth related operations.
type Service struct {
	repo repo.Repo
}

// New creates a new Service instance.
func New(repo repo.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

// Register handles logic for creating new user.
func (s *Service) Register(ctx context.Context, reqDto *dto.RegisterRequest) (*dto.User, error) {
	hashedPw, err := s.hashPassword(reqDto.Password)
	if err != nil {
		return nil, err
	}

	reqDto.Password = hashedPw

	respDto, err := s.repo.Create(ctx, reqDto)
	if err != nil {
		return nil, fmt.Errorf("creating new user: %w", err)
	}

	return respDto, nil
}

func (s *Service) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}

	return string(bytes), nil
}
