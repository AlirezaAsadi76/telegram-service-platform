package paymenthandler

import (
	"fmt"

	"github.com/labstack/echo/v5"
)

func (h *Handler) SetRoutes(e *echo.Echo) {
	fmt.Println("SetRoutes paymentHandler")
	e.GET("/payments/zarinpal/callback", h.zarinpalCallbackHandler)

}
