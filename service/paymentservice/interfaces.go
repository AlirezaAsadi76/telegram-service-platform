package paymentservice

import (
	"context"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentproviderparams"
)

type Repository interface {
	Create(ctx context.Context, payment *paymententity.Payment) error
	GetByID(ctx context.Context, id uint64) (*paymententity.Payment, error)
	GetByOrderID(ctx context.Context, orderID uint64) (*paymententity.Payment, error)
	GetByExternalID(ctx context.Context, externalID string) (*paymententity.Payment, error)
	UpdateStatus(ctx context.Context, id uint64, status paymententity.PaymentStatus) error
	GetByIdempotencyKey(ctx context.Context, key string) (*paymententity.Payment, error)
	MarkInitiated(ctx context.Context, paymentID uint64, status paymententity.PaymentStatus, externalID string, paymentURL string) error
	GetPending(ctx context.Context) ([]paymententity.Payment, error)
	GetExpired(ctx context.Context) ([]paymententity.Payment, error)
}

type Provider interface {
	Create(ctx context.Context, request paymentproviderparams.CreateRequest) (paymentproviderparams.CreateResponse, error)
	Verify(ctx context.Context, request paymentproviderparams.VerifyRequest) (paymentproviderparams.VerifyResponse, error)
}

type IdempotencyChecker interface {
	SetIfNotExists(ctx context.Context, key string, value string, ttlSeconds int) (bool, error)
}

type PaymentConfirmationRepository interface {
	Confirm(ctx context.Context, paymentID uint64, providerReferenceID string) error
	Fail(ctx context.Context, paymentID uint64) error
	MarkUnknown(ctx context.Context, paymentID uint64) error
}
