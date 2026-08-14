package redisqueue

import (
	"context"
	"encoding/json"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (db DB) LPush(ctx context.Context, queueKey string, data any) error {
	const Op = "redis_repository.LPush"

	dat, err := json.Marshal(data)
	if err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindSerializationFailure).WithMessage(msgerror.UnmarshalFailed)
	}
	if lErr := db.adapter.Client().LPush(ctx, queueKey, dat).Err(); lErr != nil {
		return richerror.New(Op, lErr).WithKind(richerror.KindCreateFailed).WithMessage(msgerror.QueryScanFailed)
	}
	return nil
}
