package redisactivity

import (
	"context"
	"strconv"
	"strings"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (t *ActivityTracker) GetActiveUsers(ctx context.Context) ([]int64, error) {
	const op = "redisactivity.GetActiveUsers"

	var (
		activeUsers []int64
		cursor      uint64
	)

	for {
		keys, nextCursor, err := t.redis.Client().Scan(ctx, cursor, activityScanPattern, 100).Result()
		if err != nil {
			return nil, richerror.New(op, err).
				WithKind(richerror.KindUnexpected).
				WithMessage(msgerror.CacheReadFailed)
		}

		for _, key := range keys {
			idStr := strings.TrimPrefix(key, activityKeyPrefix)
			telegramID, parseErr := strconv.ParseInt(idStr, 10, 64)
			if parseErr != nil {
				continue
			}
			activeUsers = append(activeUsers, telegramID)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return activeUsers, nil
}
