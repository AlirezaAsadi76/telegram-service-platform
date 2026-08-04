package redisprice

import (
	"context"
	"strconv"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (db DB) GetTonPrice(ctx context.Context) (float64, error) {
	const Op = "redisprice.GetTonPrice"
	result, aErr := db.adapter.Client().Get(ctx, TonPriceKey).Result()
	if aErr != nil {
		return 0, richerror.New(Op, aErr).
			WithKind(richerror.KindUnexpected).
			WithMessage(msgerror.CacheReadFailed)
	}
	price, err := strconv.ParseFloat(result, 64)
	if err != nil {
		return 0, richerror.New(Op, err).
			WithKind(richerror.KindUnexpected).
			WithMessage(msgerror.CacheParseFailed)
	}
	return price, nil
}
