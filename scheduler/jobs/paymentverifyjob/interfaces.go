package paymentverifyjob

import (
	"context"
)

type RedisRepository interface {
	LPush(ctx context.Context, queueKey string, data any) error
}
