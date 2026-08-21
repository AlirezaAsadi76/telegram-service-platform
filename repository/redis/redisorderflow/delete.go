package redisorderflow

import (
	"context"
	"fmt"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (db *DB) Delete(ctx context.Context, req orderparams.DeleteOrderFlowRequest) error {
	const op = "redisorderflow.Delete"

	key := fmt.Sprintf(orderFlowKeyPattern, req.TelegramID)

	if err := db.redis.Client().Del(ctx, key).Err(); err != nil {
		return richerror.New(op, err).
			WithKind(richerror.KindUnexpected).
			WithMessage(msgerror.CacheWriteFailed)
	}

	return nil
}
