package redisidempotency

import (
	"context"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d DB) Get(ctx context.Context, idempotencyKey string) (string, error) {
	const Op = "redisidempotency.Get"
	value, err := d.adapter.Client().Get(ctx, idempotencyKey).Result()
	if err != nil {
		return "", richerror.New(Op, err).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}
	return value, nil
}
