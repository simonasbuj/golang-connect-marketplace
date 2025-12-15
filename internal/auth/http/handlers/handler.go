// Package handlers defines HTTP handlers for application auth endpoints.
package handlers

import (
	"golang-connect-marketplace/internal/auth/dto"
	"golang-connect-marketplace/internal/service"
	"golang-connect-marketplace/pkg/responses"
	"golang-connect-marketplace/pkg/validation"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler handles authentication-related HTTP requests.
type Handler struct {
	svc *service.Service
}

// NewHandler creates a new Handler for handling authentication requests.
func NewHandler(svc *service.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

// HandleRegister handles requests to register user.
func (h *Handler) HandleRegister(c echo.Context) error {
	var reqDto dto.RegisterRequest

	err := validation.ValidateDto(c, &reqDto)
	if err != nil {
		return responses.JSONError(c, err.Error(), err)
	}

	respDto, err := h.svc.Register(c.Request().Context(), &reqDto)
	if err != nil {
		return responses.JSONError(
			c,
			"failed to register user",
			err,
			http.StatusInternalServerError,
		)
	}

	return responses.JSONSuccess(c, "going to register", respDto)
}

// HandleLogin handles requests to login user.
func (h *Handler) HandleLogin(c echo.Context) error {
	return responses.JSONSuccess(c, "going to login", nil)
}

// HandleRefresh handles requests to refresh token.
func (h *Handler) HandleRefresh(c echo.Context) error {
	return responses.JSONSuccess(c, "going to refresh token", nil)
}
