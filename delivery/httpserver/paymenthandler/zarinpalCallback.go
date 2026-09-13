package paymenthandler

import (
	"net/http"
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
	fieldError, vErr := h.paymentVal.ValidateZarinpalCallback(callback)

	if vErr != nil {
		msg, code := httpmsg.CodeAndMessage(vErr)
		return c.JSON(code, map[string]interface{}{

			"msg":         msg,
			"fieldErrors": fieldError,
		})

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
