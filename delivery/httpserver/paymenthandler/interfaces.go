package paymenthandler

import (
	"context"
	"telegram-service-platform/params/paymentparams"
)

type PaymentConfirmer interface {
	ConfirmPaymentByExternalID(ctx context.Context, req paymentparams.ConfirmPaymentByExternalIDRequest) (*paymentparams.ConfirmPaymentResponse, error)
}

type PaymentValidator interface {
	ValidateZarinpalCallback(req paymentparams.ZarinpalCallback) (map[string]string, error)
}
