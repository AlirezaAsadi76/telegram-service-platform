package redisprice

import (
	"context"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"
)

func (db DB) SetTonPrice(ctx context.Context, price float64, expTime time.Duration) error {
	const Op = "redisprice.SetTonPrice"
	err := db.adapter.Client().Set(ctx, TonPriceKey, price, expTime).Err()
	if err != nil {
		return richerror.New(Op, err).
			WithKind(richerror.KindUnexpected).
			WithMessage(msgerror.CacheWriteFailed)
	}
	return nil
}
