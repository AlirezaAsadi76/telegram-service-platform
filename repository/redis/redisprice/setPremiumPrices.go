package redisprice

import (
	"context"
	"encoding/json"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"
)

func (d DB) SetPremiumPrices(ctx context.Context, prices []productentity.PremiumPrice, expiration time.Duration) error {
	const Op = "redisprice.SetPremiumPrices"

	data, mErr := json.Marshal(prices)
	if mErr != nil {
		return richerror.New(Op, mErr).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}

	if err := d.adapter.Client().Set(ctx, StarPriceKey, data, expiration).Err(); err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}
	return nil
}
