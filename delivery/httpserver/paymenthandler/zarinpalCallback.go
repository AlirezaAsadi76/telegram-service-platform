package paymenthandler

import (
	"net/http"
	"telegram-service-platform/params"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/httpmsg"

	"github.com/labstack/echo/v5"
)

func (h *Handler) zarinpalCallbackHandler(c *echo.Context) error {
	var callback paymentparams.ZarinpalCallback

	if err := c.Bind(&callback); err != nil {
		msg, code := httpmsg.CodeAndMessage(err)

		return echo.NewHTTPError(code, msg)
	}

	fieldErrors, err := h.paymentVal.ValidateZarinpalCallback(
		callback,
	)
	if err != nil {
		msg, code := httpmsg.CodeAndMessage(err)

		return c.JSON(
			code,
			params.ValidationErrorResponse{
				Message:     msg,
				FieldErrors: fieldErrors,
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

	response, cErr := h.paymentService.ConfirmPaymentByExternalID(
		c.Request().Context(),
		paymentparams.ConfirmPaymentByExternalIDRequest{
			ExternalID: callback.Authority,
			CallbackData: map[string]any{
				"Status":    callback.Status,
				"Authority": callback.Authority,
			},
		},
	)
	if cErr != nil {
		return cErr
	}

	return c.JSON(http.StatusOK, response)
}
