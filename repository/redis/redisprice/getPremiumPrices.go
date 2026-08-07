package redisprice

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"

	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d DB) GetPremiumPrices(ctx context.Context) ([]productentity.PremiumPrice, error) {

	const Op = "redisprice.GetPremiumPrices"

	data, err := d.adapter.Client().Get(ctx, PremiumPriceKey).Bytes()

	if err != nil {

		if errors.Is(err, redis.Nil) {
			return nil,
				richerror.New(Op, err).
					WithKind(richerror.KindNotFound).
					WithMessage(msgerror.CacheNotFound)
		}

		return nil,
			richerror.New(Op, err).
				WithKind(richerror.KindInfrastructure).
				WithMessage(msgerror.CacheReadFailed)
	}

	var price []productentity.PremiumPrice

	err = json.Unmarshal(data, &price)
	if err != nil {

		return nil,
			richerror.New(Op, err).
				WithKind(richerror.KindInvalid).
				WithMessage(msgerror.CacheParseFailed)
	}

	return price, nil
}
