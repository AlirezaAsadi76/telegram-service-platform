package paymenthandler

import (
	"github.com/labstack/echo/v5"
)

func (h *Handler) SetRoutes(e *echo.Echo) {

	e.GET("/payments/zarinpal/callback", h.ZarinpalCallbackHandler)

}
