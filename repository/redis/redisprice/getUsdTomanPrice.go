package redisprice

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d DB) GetUsdTomanPrice(ctx context.Context) (float64, error) {

	const Op = "redisprice.GetUsdTomanPrice"

	data, err := d.adapter.Client().Get(ctx, UsdTomanPriceKey).Bytes()

	if err != nil {

		if errors.Is(err, redis.Nil) {
			return 0,
				richerror.New(Op, err).
					WithKind(richerror.KindNotFound).
					WithMessage(msgerror.CacheNotFound)
		}

		return 0,
			richerror.New(Op, err).
				WithKind(richerror.KindInfrastructure).
				WithMessage(msgerror.CacheReadFailed)
	}

	var price float64

	err = json.Unmarshal(data, &price)
	if err != nil {

		return 0,
			richerror.New(Op, err).
				WithKind(richerror.KindInvalid).
				WithMessage(msgerror.CacheParseFailed)
	}

	return price, nil
}
