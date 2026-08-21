package redisidempotency

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"

	"go.uber.org/zap"
)

func (d DB) SetIfNotExists(ctx context.Context, idempotencyKey string, Value entity.IdempotencyStatus, ttl time.Duration) (bool, error) {
	const Op = "redisidempotency.SetIfNotExists"

	// ✅ اضافه کردن لاگ برای debug
	logger.Logger.Debug("SetIfNotExists called",
		zap.String("op", Op),
		zap.String("key", idempotencyKey[:16]+"..."),
		zap.Any("value", Value),
		zap.Duration("ttl", ttl),
	)

	exist, err := d.adapter.Client().SetNX(ctx, idempotencyKey, Value.String(), ttl).Result()
	if err != nil {
		logger.Logger.Error("SetIfNotExists failed",
			zap.String("op", Op),
			zap.String("key", idempotencyKey[:16]+"..."),
			zap.Error(err),
		)
		return false, richerror.New(Op, err).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}

	// ✅ لاگ نتیجه
	logger.Logger.Debug("SetIfNotExists result",
		zap.String("op", Op),
		zap.String("key", idempotencyKey[:16]+"..."),
		zap.Bool("created", exist),
	)

	return exist, nil
}
