// Package handlers defines HTTP handlers for application auth endpoints.
package handlers

import (
	"errors"
	"golang-connect-marketplace/internal/auth/dto"
	"golang-connect-marketplace/internal/auth/middleware"
	"golang-connect-marketplace/internal/auth/service"
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

	return responses.JSONSuccess(c, "new user registered", respDto)
}

// HandleLogin handles requests to login user.
func (h *Handler) HandleLogin(c echo.Context) error {
	var reqDto dto.LoginRequest

	err := validation.ValidateDto(c, &reqDto)
	if err != nil {
		return responses.JSONError(c, err.Error(), err)
	}

	respDto, err := h.svc.Login(c.Request().Context(), &reqDto)
	if err != nil {
		return responses.JSONError(
			c,
			"failed to login user",
			err,
			http.StatusInternalServerError,
		)
	}

	return responses.JSONSuccess(c, "user logged in successfully", respDto)
}

// HandleRefresh handles requests to refresh token.
func (h *Handler) HandleRefresh(c echo.Context) error {
	var reqDto dto.RefreshTokenRequest

	err := validation.ValidateDto(c, &reqDto)
	if err != nil {
		return responses.JSONError(c, err.Error(), err)
	}

	respDto, err := h.svc.RefreshToken(c.Request().Context(), &reqDto)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			return responses.JSONError(c, "unauthorized", err)
		}

		return responses.JSONError(
			c,
			"failed to refresh token",
			err,
			http.StatusInternalServerError,
		)
	}

	return responses.JSONSuccess(c, "refreshed tokens successfully", respDto)
}

// HandleSecret handles auth middleware debugging.
func (h *Handler) HandleSecret(c echo.Context) error {
	userClaims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	return responses.JSONSuccess(c, "can access", userClaims)
}
