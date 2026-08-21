package orderflowservice

import (
	"context"
	"telegram-service-platform/params/orderparams"
	"time"

	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/richerror"

	"telegram-service-platform/logger"

	"go.uber.org/zap"
)

func (s *Service) DeleteOrderFlow(ctx context.Context, req orderparams.DeleteOrderFlowRequest, reason string) error {
	const op = "orderflowservice.DeleteOrderFlow"
	start := time.Now()

	if req.TelegramID <= 0 {
		return richerror.New(op, nil).
			WithKind(richerror.KindValidation)
	}

	if err := s.repo.Delete(ctx, req); err != nil {
		metrics.OrderFlowStateDeleted.WithLabelValues(reason, "error").Inc()
		logger.Logger.Error("failed to delete order flow state",
			zap.String("op", op),
			zap.Int64("telegram_id", req.TelegramID.Int64()),
			zap.String("reason", reason),
			zap.Error(err),
		)
		return richerror.New(op, err)
	}

	metrics.OrderFlowStateDeleted.WithLabelValues(reason, "success").Inc()
	logger.Logger.Info("order flow state deleted",
		zap.String("op", op),
		zap.Int64("telegram_id", req.TelegramID.Int64()),
		zap.String("reason", reason),
		zap.Duration("latency", time.Since(start)),
	)

	return nil
}
