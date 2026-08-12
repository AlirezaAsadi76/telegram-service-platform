package redisidempotency

import (
	"context"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"
)

func (d *DB) SetIfNotExists(ctx context.Context, idempotencyKey string, Value string, ttl time.Duration) error {
	const Op = "redisidempotency.Set"

	if err := d.adapter.Client().SetNX(ctx, idempotencyKey, Value, ttl).Err(); err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}
	return nil
}
