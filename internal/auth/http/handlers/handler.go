// Package handlers defines HTTP handlers for application auth endpoints.
package handlers

import (
	"errors"
	"golang-connect-marketplace/config"
	"golang-connect-marketplace/internal/auth/dto"
	"golang-connect-marketplace/internal/auth/middleware"
	"golang-connect-marketplace/internal/auth/service"
	r "golang-connect-marketplace/pkg/responses"
	"golang-connect-marketplace/pkg/validation"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const refreshTokenCookieName = "refresh_token"

// ErrURLMissingParam is returned when url is missing required param.
var ErrURLMissingParam = errors.New("url is missing required param")

// Handler handles authentication-related HTTP requests.
type Handler struct {
	svc *service.Service
	cfg *config.AuthConfig
}

// NewHandler creates a new Handler for handling authentication requests.
func NewHandler(svc *service.Service, cfg *config.AuthConfig) *Handler {
	return &Handler{
		svc: svc,
		cfg: cfg,
	}
}

// HandleRegister handles requests to register user.
func (h *Handler) HandleRegister(c echo.Context) error {
	var reqDto dto.RegisterRequest

	err := validation.ValidateDto(c, &reqDto)
	if err != nil {
		return r.JSONError(c, err.Error(), err)
	}

	respDto, err := h.svc.Register(c.Request().Context(), &reqDto)
	if err != nil {
		return r.JSONError(c, "failed to register user", err, http.StatusInternalServerError)
	}

	return r.JSONSuccess(c, "new user registered", respDto)
}

// HandleLogin handles requests to login user.
func (h *Handler) HandleLogin(c echo.Context) error {
	var reqDto dto.LoginRequest

	err := validation.ValidateDto(c, &reqDto)
	if err != nil {
		return r.JSONError(c, err.Error(), err)
	}

	respDto, err := h.svc.Login(c.Request().Context(), &reqDto)
	if err != nil {
		return r.JSONError(c, "failed to login user", err, http.StatusInternalServerError)
	}

	refreshTokenCookie := h.createTokenCookie(refreshTokenCookieName, respDto.RefreshToken)
	c.SetCookie(refreshTokenCookie)

	return r.JSONSuccess(c, "user logged in successfully", respDto)
}

// HandleRefresh handles requests to refresh token.
func (h *Handler) HandleRefresh(c echo.Context) error {
	refreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil {
		return r.JSONError(c, "refresh token missing in cookies", err, http.StatusUnauthorized)
	}

	reqDto := &dto.RefreshTokenRequest{
		RefreshToken: refreshToken.Value,
	}

	respDto, err := h.svc.RefreshToken(c.Request().Context(), reqDto)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			return r.JSONError(c, "unauthorized", err, http.StatusUnauthorized)
		}

		return r.JSONError(c, "failed to refresh token", err, http.StatusInternalServerError)
	}

	refreshTokenCookie := h.createTokenCookie(refreshTokenCookieName, respDto.RefreshToken)
	c.SetCookie(refreshTokenCookie)

	return r.JSONSuccess(c, "refreshed tokens successfully", respDto)
}

// HandleSecret handles auth middleware debugging.
func (h *Handler) HandleSecret(c echo.Context) error {
	userClaims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	return r.JSONSuccess(c, "can access", userClaims)
}

func (h *Handler) createTokenCookie(cookieName, token string) *http.Cookie {
	return &http.Cookie{ //nolint:exhaustruct
		Name:     cookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   h.cfg.RefreshTokenCookieSecure,
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(h.cfg.RefreshTokenValidSeconds) * time.Second),
		MaxAge:   h.cfg.RefreshTokenValidSeconds,
	}
}
