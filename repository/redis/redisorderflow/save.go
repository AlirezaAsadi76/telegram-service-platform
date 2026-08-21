package redisorderflow

import (
	"context"
	"encoding/json"
	"fmt"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (db *DB) Save(ctx context.Context, req orderparams.SaveOrderFlowRequest) error {
	const op = "redisorderflow.Save"

	key := fmt.Sprintf(orderFlowKeyPattern, req.TelegramID)

	data, err := json.Marshal(req.State)
	if err != nil {
		return richerror.New(op, err).
			WithKind(richerror.KindSerializationFailure).
			WithMessage(msgerror.MarshalFailed)
	}

	if err := db.redis.Client().Set(ctx, key, data, req.TTLMins).Err(); err != nil {
		return richerror.New(op, err).
			WithKind(richerror.KindUnexpected).
			WithMessage(msgerror.CacheWriteFailed)
	}

	return nil
}
