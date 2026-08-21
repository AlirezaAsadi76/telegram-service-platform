package checkoutservice

import (
	"context"
	"telegram-service-platform/entity"
	"time"

	"github.com/go-telegram/bot"
)

type Messenger interface {
	Send(ctx context.Context, params *bot.SendMessageParams) error
}

type IdempotencyChecker interface {
	SetIfNotExists(ctx context.Context, key string, Value entity.IdempotencyStatus, ttl time.Duration) (bool, error)
	Set(ctx context.Context, key string, Value entity.IdempotencyStatus, ttl time.Duration) error
}
