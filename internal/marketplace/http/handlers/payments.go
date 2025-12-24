package handlers

import (
	"golang-connect-marketplace/internal/auth/middleware"
	"golang-connect-marketplace/internal/marketplace/dto"
	"golang-connect-marketplace/internal/marketplace/services"
	"golang-connect-marketplace/pkg/validation"
	"net/http"

	r "golang-connect-marketplace/pkg/responses"

	"github.com/labstack/echo/v4"
)

// PaymentsHandler handles payments-related HTTP requests.
type PaymentsHandler struct {
	svc *services.PaymentsService
}

// NewPaymentsHandler creates a new Handler for handling payments requests.
func NewPaymentsHandler(svc *services.PaymentsService) *PaymentsHandler {
	return &PaymentsHandler{
		svc: svc,
	}
}

// HandleLinkSellerAccount handles linking user to seller account.
func (h *PaymentsHandler) HandleLinkSellerAccount(c echo.Context) error {
	userClaims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	var reqDto dto.SellerAcountLinkingSessionRequest

	reqDto.UserID = userClaims.ID

	err = validation.ValidateDto(c, &reqDto)
	if err != nil {
		return r.JSONError(c, err.Error(), err)
	}

	resp, err := h.svc.LinkSellerAccount(c.Request().Context(), &reqDto)
	if err != nil {
		return r.JSONError(
			c,
			"failed to create linking session",
			err,
			http.StatusInternalServerError,
		)
	}

	return r.JSONSuccess(c, "created linking session", resp)
}
