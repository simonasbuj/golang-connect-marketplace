// Package middleware provides HTTP middleware.
package middleware

import (
	"errors"
	"golang-connect-marketplace/internal/auth/dto"
	"golang-connect-marketplace/internal/auth/service"
	"golang-connect-marketplace/pkg/responses"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// ErrUserClaimsNotFoundInContext is returned when user claims are not found in echo context.
var ErrUserClaimsNotFoundInContext = errors.New("user claims not found in context")

const userContextKey = "user"

// AuthenticateMiddleware checks if token provided in Authorization header is valid
// and saves UserClaims in echo context.
func AuthenticateMiddleware(svc *service.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return responses.JSONError(
					c,
					"no token in header: Authorization",
					service.ErrUnauthorized,
					http.StatusUnauthorized,
				)
			}

			token := strings.TrimPrefix(auth, "Bearer ")

			userClaims, err := svc.ParseJWT(token)
			if err != nil {
				return responses.JSONError(
					c,
					"invalid token",
					err,
					http.StatusUnauthorized,
				)
			}

			c.Set(userContextKey, userClaims)

			return next(c)
		}
	}
}

// GetUserFromContext fetches UserClaims object from echo context.
func GetUserFromContext(c echo.Context) (*dto.UserClaims, error) {
	claims, ok := c.Get(userContextKey).(*dto.UserClaims)
	if !ok || claims == nil {
		return nil, responses.JSONError( //nolint:wrapcheck
			c,
			"unauthorized",
			ErrUserClaimsNotFoundInContext,
			http.StatusUnauthorized,
		)
	}

	return claims, nil
}
