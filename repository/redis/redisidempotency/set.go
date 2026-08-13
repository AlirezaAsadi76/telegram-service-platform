package redisidempotency

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"
)

func (d DB) Set(ctx context.Context, idempotencyKey string, Value entity.IdempotencyStatus, ttl time.Duration) error {
	const Op = "redisidempotency.Set"

	if err := d.adapter.Client().Set(ctx, idempotencyKey, Value, ttl).Err(); err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}
	return nil
}
