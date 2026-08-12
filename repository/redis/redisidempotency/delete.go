package redisidempotency

import (
	"context"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) Delete(ctx context.Context, idempotencyKey string) error {
	const Op = "redisidempotency.Delete"

	if err := d.adapter.Client().Del(ctx, idempotencyKey).Err(); err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}
	return nil
}
