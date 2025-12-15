package handlers

import (
	"golang-connect-marketplace/pkg/responses"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) HandleRegister(c echo.Context) error {
	return responses.JSONSuccess(c, "going to register", nil)
}

func (h *AuthHandler) HandleLogin(c echo.Context) error {
	return responses.JSONSuccess(c, "going to login", nil)
}

func (h *AuthHandler) HandleRefresh(c echo.Context) error {
	return responses.JSONSuccess(c, "going to refresh token", nil)
}
