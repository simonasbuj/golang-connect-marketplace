// Package service implements the core business logic of authentication and user management.
package service

import (
	"context"
	"fmt"
	"golang-connect-marketplace/internal/auth/dto"
	"golang-connect-marketplace/pkg/generate"

	"golang.org/x/crypto/bcrypt"
)

// Service provides user and auth related operations.
type Service struct{}

// NewService creates a new Service instance.
func NewService() *Service {
	return &Service{}
}

// Register handles logic for creating new user.
func (s *Service) Register(_ context.Context, reqDto *dto.RegisterRequest) (*dto.User, error) {
	hashedPw, err := s.hashPassword(reqDto.Password)
	if err != nil {
		return nil, err
	}

	reqDto.Password = hashedPw

	respDto := &dto.User{
		ID:       generate.ID("user"),
		Email:    reqDto.Email,
		Name:     reqDto.Name,
		Lastname: reqDto.Lastname,
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
