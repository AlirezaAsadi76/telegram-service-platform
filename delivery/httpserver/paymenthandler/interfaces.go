package paymenthandler

import (
	"context"
	"telegram-service-platform/params/paymentparams"

	"github.com/labstack/echo/v5"
)

type PaymentFlow interface {
	StartPayment(ctx context.Context, req paymentparams.StartPaymentRequest) (*paymentparams.InitiateResponse, error)
	ConfirmPaymentByExternalID(ctx context.Context, req paymentparams.ConfirmPaymentByExternalIDRequest) (*paymentparams.ConfirmPaymentResponse, error)
}

type PaymentValidator interface {
	ValidateZarinpalCallback(req paymentparams.ZarinpalCallback) (map[string]string, error)
	ValidateStartPayment(req paymentparams.StartPaymentHandlerRequest) (map[string]string, error)
}

type AuthMiddleware interface {
	Auth() echo.MiddlewareFunc
}
