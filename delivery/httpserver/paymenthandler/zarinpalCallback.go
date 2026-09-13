package paymenthandler

import (
	"net/http"
	"telegram-service-platform/params/paymentparams"

	"github.com/labstack/echo/v5"
)

func (h *Handler) zarinpalCallbackHandler(c *echo.Context) error {

	var callback paymentparams.ZarinpalCallback

	if callback.Authority == "" {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error":   "invalid_callback",
				"message": "authority is required",
			},
		)
	}

	if callback.Status != "OK" {
		return c.JSON(
			http.StatusOK,
			map[string]string{
				"status": "cancelled",
			},
		)
	}

	response, err := h.paymentService.ConfirmPaymentByExternalID(
		c.Request().Context(),
		paymentparams.ConfirmPaymentByExternalIDRequest{
			ExternalID: callback.Authority,
			CallbackData: map[string]any{
				"Status":    callback.Status,
				"Authority": callback.Authority,
			},
		},
	)
	if err != nil {
		return err
	}

	return c.JSON(
		http.StatusOK,
		response,
	)
}
