package paymenthandler

import (
	"context"
	"telegram-service-platform/params/paymentparams"

	"github.com/labstack/echo/v5"
)

type PaymentConfirmer interface {
	ConfirmPaymentByExternalID(ctx context.Context, req paymentparams.ConfirmPaymentByExternalIDRequest) (*paymentparams.ConfirmPaymentResponse, error)
}

type Handler struct {
	paymentService PaymentConfirmer
}

func New(
	paymentService PaymentConfirmer,
) *Handler {
	return &Handler{
		paymentService: paymentService,
	}
}
