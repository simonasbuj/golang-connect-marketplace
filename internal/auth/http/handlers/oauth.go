package handlers

import (
	"golang-connect-marketplace/internal/auth/dto"
	r "golang-connect-marketplace/pkg/responses"
	"golang-connect-marketplace/pkg/validation"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HandleGithubInit handles initial GitHub oauth request.
func (h *Handler) HandleGithubInit(c echo.Context) error {
	url := h.svc.GetGithubAuthURL()

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleOauthCallback handles oauth callback for all oatuh providers.
func (h *Handler) HandleOauthCallback(c echo.Context) error {
	var reqDto dto.OauthExchangeRequest

	err := validation.ValidateDto(c, &reqDto)
	if err != nil {
		return r.JSONError(c, "missing exchange code", err)
	}

	return r.JSONSuccess(c, "Oauth callback handled", reqDto.Code)
}

// HandleGithubExchange handles GitHub oauth code exchange for jwt token request.
func (h *Handler) HandleGithubExchange(c echo.Context) error {
	var reqDto dto.OauthExchangeRequest

	err := validation.ValidateDto(c, &reqDto)
	if err != nil {
		return r.JSONError(c, "missing exchange code", err)
	}

	resp, err := h.svc.HandleGithubCallback(c.Request().Context(), reqDto.Code)
	if err != nil {
		return r.JSONError(c, "failed to handle github oauth exchange", err)
	}

	refreshTokenCookie := h.createRefreshTokenCookie(resp.RefreshToken)
	c.SetCookie(refreshTokenCookie)

	return r.JSONSuccess(c, "github oauth exchange handled", resp)
}

// HandleGoogleInit handles initial Google oauth request.
func (h *Handler) HandleGoogleInit(c echo.Context) error {
	url := h.svc.GetGoogleAuthURL()

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleGoogleExchange handles Google oauth code exchange for jwt token request.
func (h *Handler) HandleGoogleExchange(c echo.Context) error {
	var reqDto dto.OauthExchangeRequest

	err := validation.ValidateDto(c, &reqDto)
	if err != nil {
		return r.JSONError(c, "missing exchange code", err)
	}

	resp, err := h.svc.HandleGoogleCallback(c.Request().Context(), reqDto.Code)
	if err != nil {
		return r.JSONError(c, "failed to handle google oauth exchange", err)
	}

	refreshTokenCookie := h.createRefreshTokenCookie(resp.RefreshToken)
	c.SetCookie(refreshTokenCookie)

	return r.JSONSuccess(c, "google oauth exchange handled", resp)
}
