package redisactivity

import (
	"context"
	"fmt"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (t *ActivityTracker) TrackActivity(ctx context.Context, telegramID int64) error {
	const op = "redisactivity.TrackActivity"

	key := fmt.Sprintf(activityKeyPattern, telegramID)

	if err := t.redis.Client().Set(ctx, key, "online", t.ttl).Err(); err != nil {
		return richerror.New(op, err).
			WithKind(richerror.KindUnexpected).
			WithMessage(msgerror.CacheWriteFailed)
	}

	return nil
}
