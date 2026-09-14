package paymenthandler

import (
	"net/http"
	"telegram-service-platform/params"
	"telegram-service-platform/params/paymentparams"
	claimspkg "telegram-service-platform/pkg/claims"
	"telegram-service-platform/pkg/httpmsg"

	"github.com/labstack/echo/v5"
)

func (h *Handler) StartPaymentHandler(c *echo.Context) error {
	var req paymentparams.StartPaymentHandlerRequest

	if err := c.Bind(&req); err != nil {
		msg, code := httpmsg.CodeAndMessage(err)

		return echo.NewHTTPError(code, msg)
	}

	fieldErrors, err := h.paymentVal.ValidateStartPayment(req)
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

	claims := claimspkg.GetClaimsFromEchoContext(c)

	if claims == nil {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"unauthorized",
		)
	}

	response, sErr := h.paymentService.StartPayment(
		c.Request().Context(),
		paymentparams.StartPaymentRequest{
			OrderID:        req.OrderID,
			UserID:         uint64(claims.UserId),
			Method:         req.Method,
			Amount:         req.Amount,
			Currency:       req.Currency,
			IdempotencyKey: req.IdempotencyKey,
			CallbackURL:    req.CallbackURL,
			Description:    req.Description,
		},
	)
	if sErr != nil {
		msg, code := httpmsg.CodeAndMessage(sErr)

		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, response)
}
