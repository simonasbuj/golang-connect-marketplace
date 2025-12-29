package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"golang-connect-marketplace/internal/auth/dto"

	// "golang-connect-marketplace/internal/auth/dto"
	"net/http"
	"net/url"
	"strings"
)

func (s *Service) HandleGithubCallback(ctx context.Context, code string) (*dto.User, error) {
	token, err := s.exchangeCodeForToken(code)
	if err != nil {
		return nil, fmt.Errorf("exchanging github callback code for token: %w", err)
	}

	ghUser, err := s.fetchGithubUser(token)
	if err != nil {
		return nil, fmt.Errorf("fetching github user: %w", err)
	}

	user, err := s.repo.GetUserByOAuthID(ctx, string(ghUser.ID), dto.OAuthProviderGitHub)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("fetching user by oauth id: %w", err)
	}

	return user, nil
}

type GithubTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *Service) exchangeCodeForToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", s.cfg.OAuthConfig.GithubClientID)
	data.Set("client_secret", s.cfg.OAuthConfig.GithubClientSecret)
	data.Set("code", code)

	req, err := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res GithubTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.AccessToken, nil
}

type GithubUser struct {
	ID    int  `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}

func (s *Service) fetchGithubUser(token string) (*GithubUser, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user GithubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}