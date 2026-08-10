package paymentservice

import (
	"context"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
)

type Repository interface {
	Create(ctx context.Context, payment *paymententity.Payment) error
	GetByID(ctx context.Context, id uint64) (*paymententity.Payment, error)
	GetByOrderID(ctx context.Context, orderID uint64) (*paymententity.Payment, error)
	UpdateStatus(ctx context.Context, id uint64, status paymententity.PaymentStatus) error
}

type Provider interface {
	Create(ctx context.Context, request paymentparams.CreatePaymentRequest) (paymentparams.CreatePaymentResponse, error)
	Verify(ctx context.Context, request paymentparams.VerifyPaymentRequest) (paymentparams.VerifyPaymentResponse, error)
}
