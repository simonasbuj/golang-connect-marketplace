// Package handlers defines HTTP handlers for application auth endpoints.
package handlers

import (
	"golang-connect-marketplace/pkg/responses"

	"github.com/labstack/echo/v4"
)

// Handler handles authentication-related HTTP requests.
type Handler struct{}

// NewHandler creates a new Handler for handling authentication requests.
func NewHandler() *Handler {
	return &Handler{}
}

// HandleRegister handles requests to register user.
func (h *Handler) HandleRegister(c echo.Context) error {
	return responses.JSONSuccess(c, "going to register", nil)
}

// HandleLogin handles requests to login user.
func (h *Handler) HandleLogin(c echo.Context) error {
	return responses.JSONSuccess(c, "going to login", nil)
}

// HandleRefresh handles requests to refresh token.
func (h *Handler) HandleRefresh(c echo.Context) error {
	return responses.JSONSuccess(c, "going to refresh token", nil)
}
