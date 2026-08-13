package checkoutservice

import (
	"context"
	"telegram-service-platform/entity"
	"time"
)

type Messenger interface {
	SendToUser(ctx context.Context, userID int64, message string) error
	SendToAdminChannel(ctx context.Context, message string) error
}

type IdempotencyChecker interface {
	SetIfNotExists(ctx context.Context, key string, Value entity.IdempotencyStatus, ttl time.Duration) (bool, error)
	Set(ctx context.Context, key string, Value entity.IdempotencyStatus, ttl time.Duration) error
}
