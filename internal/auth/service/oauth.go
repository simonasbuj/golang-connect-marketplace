package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"golang-connect-marketplace/internal/auth/dto"
	"net/http"
	"strconv"
	"strings"

	// "net/http".

	"golang.org/x/oauth2"
)

// ErrOAuthAPIBadStatus is returned when external oauth api returns status other than OK (200).
var ErrOAuthAPIBadStatus = errors.New("external oauth API returned bad status")

// GetGithubAuthURL fetches endpoint for github authorization.
func (s *Service) GetGithubAuthURL() string {
	url := s.oauthGithub.AuthCodeURL("", oauth2.AccessTypeOffline)

	return url
}

// HandleGithubCallback handles logic for github oatuh callback.
func (s *Service) HandleGithubCallback(
	ctx context.Context,
	code string,
) (*dto.LoginResponse, error) {
	token, err := s.oauthGithub.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchanging github callback code for token: %w", err)
	}

	ghUser, err := fetchExternalOAuthUser[dto.GithubUser](
		ctx,
		token,
		"https://api.github.com/user",
		"github",
	)
	if err != nil {
		return nil, fmt.Errorf("fetching github user: %w", err)
	}

	parts := strings.SplitN(ghUser.Name, " ", 2) //nolint:mnd
	firstName := parts[0]

	lastName := "Unknown"
	if len(parts) > 1 && strings.TrimSpace(parts[1]) != "" {
		lastName = parts[1]
	}

	user, err := s.getOrCreateOAuthUser(
		ctx,
		strconv.Itoa(ghUser.ID),
		ghUser.Email,
		firstName,
		lastName,
		ghUser.Login,
		dto.OAuthProviderGitHub,
	)
	if err != nil {
		return nil, fmt.Errorf("getting or creating user for github oauth: %w", err)
	}

	respDto, err := s.createTokens(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("creting tokens: %w", err)
	}

	return respDto, nil
}

// GetGoogleAuthURL fetches endpoint for google authorization.
func (s *Service) GetGoogleAuthURL() string {
	url := s.oauthGoogle.AuthCodeURL("", oauth2.AccessTypeOffline)

	return url
}

// HandleGoogleCallback handles logic for google oauth callback.
func (s *Service) HandleGoogleCallback(
	ctx context.Context,
	code string,
) (*dto.LoginResponse, error) {
	token, err := s.oauthGoogle.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchanging github callback code for token: %w", err)
	}

	googleUser, err := fetchExternalOAuthUser[dto.GoogleUser](
		ctx,
		token,
		"https://www.googleapis.com/oauth2/v2/userinfo",
		"google",
	)
	if err != nil {
		return nil, fmt.Errorf("fetching google user: %w", err)
	}

	user, err := s.getOrCreateOAuthUser(
		ctx,
		googleUser.ID,
		googleUser.Email,
		googleUser.GivenName,
		googleUser.FamilyName,
		googleUser.Name,
		dto.OAuthProviderGoogle,
	)
	if err != nil {
		return nil, fmt.Errorf("getting or creating user for github oauth: %w", err)
	}

	resp, err := s.createTokens(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("creting tokens: %w", err)
	}

	return resp, nil
}

func (s *Service) getOrCreateOAuthUser(
	ctx context.Context,
	oauthID, email, name, lastname, username string,
	provider dto.OAuthProvider,
) (*dto.User, error) {
	user, err := s.repo.GetUserByOAuthID(ctx, oauthID, provider)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("fetching user by oauth id: %w", err)
	}

	if user != nil {
		return user, nil
	}

	user, err = s.repo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("fetching user by email: %w", err)
	}

	if user != nil {
		_, err = s.repo.CreateOAuthUser(ctx, &dto.OAuthUser{ //nolint:exhaustruct
			UserID:         user.ID,
			ProviderUserID: oauthID,
			Provider:       provider,
		})
		if err != nil {
			return nil, fmt.Errorf("creating new oauth user: %w", err)
		}

		return user, nil
	}

	user, err = s.repo.CreateUser(ctx, &dto.RegisterRequest{ //nolint:exhaustruct
		Email:    email,
		Password: "_",
		Name:     name,
		Lastname: lastname,
		Username: username,
	})
	if err != nil {
		return nil, fmt.Errorf("creating new user during oauth: %w", err)
	}

	_, err = s.repo.CreateOAuthUser(ctx, &dto.OAuthUser{ //nolint:exhaustruct
		UserID:         user.ID,
		ProviderUserID: oauthID,
		Provider:       provider,
	})
	if err != nil {
		return nil, fmt.Errorf("creating new oauth user: %w", err)
	}

	return user, nil
}

func fetchExternalOAuthUser[T any](
	ctx context.Context,
	token *oauth2.Token,
	url string,
	providerName string,
) (*T, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating %s request: %w", providerName, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching %s user: %w", providerName, err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrOAuthAPIBadStatus, resp.StatusCode)
	}

	var resource T

	err = json.NewDecoder(resp.Body).Decode(&resource)
	if err != nil {
		return nil, fmt.Errorf("decoding %s response: %w", providerName, err)
	}

	return &resource, nil
}
