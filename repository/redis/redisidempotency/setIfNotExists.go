package redisidempotency

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"
)

func (d DB) SetIfNotExists(ctx context.Context, idempotencyKey string, Value entity.IdempotencyStatus, ttl time.Duration) (bool, error) {
	const Op = "redisidempotency.Set"

	exist, err := d.adapter.Client().SetNX(ctx, idempotencyKey, Value, ttl).Result()

	if err != nil {
		return false, richerror.New(Op, err).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}
	return exist, nil
}
