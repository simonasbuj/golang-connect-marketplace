// Package service implements the core business logic of authentication and user management.
package service

import (
	"context"
	"fmt"
	"golang-connect-marketplace/config"
	"golang-connect-marketplace/internal/auth/dto"
	"golang-connect-marketplace/internal/auth/repo"
	"golang-connect-marketplace/pkg/generate"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Service provides user and auth related operations.
type Service struct {
	repo repo.Repo
	cfg  *config.AuthConfig
}

// New creates a new Service instance.
func New(repo repo.Repo, cfg *config.AuthConfig) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}

// Register handles logic for creating new user.
func (s *Service) Register(ctx context.Context, reqDto *dto.RegisterRequest) (*dto.User, error) {
	hashedPw, err := s.hashPassword(reqDto.Password)
	if err != nil {
		return nil, err
	}

	reqDto.Password = hashedPw

	respDto, err := s.repo.CreateUser(ctx, reqDto)
	if err != nil {
		return nil, fmt.Errorf("creating new user: %w", err)
	}

	return respDto, nil
}

// Login handles logic for logging in the user.
func (s *Service) Login(ctx context.Context, reqDto *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, reqDto.Email)
	if err != nil {
		return nil, fmt.Errorf("fetching user: %w", err)
	}

	err = s.verifyPassword(user.PasswordHash, reqDto.Password)
	if err != nil {
		return nil, fmt.Errorf("verifying password: %w", err)
	}

	accessToken, err := s.generateJWT(user, s.cfg.TokenValidSeconds)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken := &dto.RefreshToken{
		Token:     strings.Trim(generate.ID("", s.cfg.RefreshTokenLength), "_"),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.RefreshTokenValidSeconds) * time.Second),
	}

	err = s.repo.SaveRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("saving refresh token: %w", err)
	}

	respDto := &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
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

func (s *Service) verifyPassword(passwordHash, passwordStr string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordStr))
	if err != nil {
		return fmt.Errorf("comparing hash and password: %w", err)
	}

	return nil
}

func (s *Service) generateJWT(user *dto.User, ttlSeconds int) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Duration(ttlSeconds) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtToken, err := token.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", fmt.Errorf("signing jtw token string: %w", err)
	}

	return jwtToken, nil
}
