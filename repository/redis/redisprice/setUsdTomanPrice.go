package redisprice

import (
	"context"
	"encoding/json"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"
)

func (d DB) SetUsdTomanPrice(ctx context.Context, price float64, expiration time.Duration) error {
	const Op = "redisprice.SetUsdTomanPrice"

	data, mErr := json.Marshal(price)

	if mErr != nil {
		return richerror.New(Op, mErr).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}

	if err := d.adapter.Client().Set(ctx, UsdTomanPriceKey, data, expiration).Err(); err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}
	return nil
}
