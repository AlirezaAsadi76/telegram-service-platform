package redisprice

import (
	"context"
	"encoding/json"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"
)

func (a DB) SetStarPrice(ctx context.Context, price productentity.StarPrice, expiration time.Duration) error {
	const Op = "redisprice.SetStarPrice"

	data, mErr := json.Marshal(price)
	if mErr != nil {
		return richerror.New(Op, mErr).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}

	if err := a.adapter.Client().Set(ctx, StarPriceKey, data, expiration).Err(); err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}
	return nil
}
