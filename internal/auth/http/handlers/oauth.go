package handlers

import (
	"fmt"
	r "golang-connect-marketplace/pkg/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HandleGithub handles GitHub oauth request.
func (h *Handler) HandleGithub(c echo.Context) error {
	url := h.svc.GetGithubAuthURL()

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleGithubCallback handles GitHub oauth callback request.
func (h *Handler) HandleGithubCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return r.JSONError(
			c,
			"missing code param in callback url",
			fmt.Errorf("%w: code: %s", ErrURLMissingParam, "code"),
		)
	}

	resp, err := h.svc.HandleGithubCallback(c.Request().Context(), code)
	if err != nil {
		return r.JSONError(c, "failed to handle github oauth callback", err)
	}

	return r.JSONSuccess(c, "github oauth callback handled", resp)
}

// HandleGoogle handles Google oauth request.
func (h *Handler) HandleGoogle(c echo.Context) error {
	url := h.svc.GetGoogleAuthURL()

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleGoogleCallback handles Google oauth callback request.
func (h *Handler) HandleGoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return r.JSONError(
			c,
			"missing code param in callback url",
			fmt.Errorf("%w: code: %s", ErrURLMissingParam, "code"),
		)
	}

	resp, err := h.svc.HandleGoogleCallback(c.Request().Context(), code)
	if err != nil {
		return r.JSONError(c, "failed to handle google oauth callback", err)
	}

	return r.JSONSuccess(c, "google oauth callback handled", resp)
}
