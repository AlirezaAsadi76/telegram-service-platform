package checkoutservice

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/params/walletparam"
	"time"
)

type IdempotencyChecker interface {
	SetIfNotExists(ctx context.Context, key string, Value entity.IdempotencyStatus, ttl time.Duration) (bool, error)
	Set(ctx context.Context, key string, Value entity.IdempotencyStatus, ttl time.Duration) error
}

type WalletPurchaseRepository interface {
	ExecuteWalletPurchase(ctx context.Context, req walletparam.WalletPurchaseRequest) (*walletparam.WalletPurchaseResult, error)
}

type FulfillmentEnqueuer interface {
	Enqueue(ctx context.Context, orderID uint64) error
}
