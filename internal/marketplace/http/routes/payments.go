package routes

import (
	m "golang-connect-marketplace/internal/auth/middleware"
	"golang-connect-marketplace/internal/auth/service"
	"golang-connect-marketplace/internal/marketplace/http/handlers"

	"github.com/labstack/echo/v4"
)

// RegisterPaymentsRoutes registers marketplace-related HTTP routes.
func RegisterPaymentsRoutes(e *echo.Echo, h *handlers.PaymentsHandler, authSvc *service.Service) {
	api := e.Group("api/v1/payments")

	api.POST("/link-seller", h.HandleLinkSellerAccount, m.AuthenticateMiddleware(authSvc))
	api.POST("/:listing_id", h.HandleCreateCheckoutSession, m.AuthenticateMiddleware(authSvc))

	api.POST("/webhook/success", h.HandlePaymentWebhookSuccess)
	api.POST("/webhook/refund", h.HandlePaymentWebhookRefund)
}
