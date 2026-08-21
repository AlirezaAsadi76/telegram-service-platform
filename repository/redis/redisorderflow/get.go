package redisorderflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/redis/go-redis/v9"
)

func (db *DB) Get(ctx context.Context, req orderparams.GetOrderFlowRequest) (*orderentity.OrderFlowState, error) {
	const op = "redisorderflow.Get"

	key := fmt.Sprintf(orderFlowKeyPattern, req.TelegramID)

	data, err := db.redis.Client().Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // NotFound
		}
		return nil, richerror.New(op, err).
			WithKind(richerror.KindUnexpected).
			WithMessage(msgerror.CacheReadFailed)
	}

	var state orderentity.OrderFlowState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, richerror.New(op, err).
			WithKind(richerror.KindSerializationFailure).
			WithMessage(msgerror.UnmarshalFailed)
	}

	return &state, nil
}
