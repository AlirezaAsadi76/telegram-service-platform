package orderflowservice

import (
	"context"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/metrics"
	"time"

	"go.uber.org/zap"
)

func (s *Service) CompleteOrderFlow(ctx context.Context, req orderparams.DeleteOrderFlowRequest, duration time.Duration) error {
	const op = "orderflowservice.CompleteOrderFlow"

	metrics.OrderFlowDuration.WithLabelValues("completed").Observe(duration.Seconds())

	if err := s.DeleteOrderFlow(ctx, req, "completed"); err != nil {
		logger.Logger.Error("failed to complete order flow",
			zap.String("op", op),
			zap.Int64("telegram_id", req.TelegramID.Int64()),
			zap.Error(err),
		)
		return err
	}

	logger.Logger.Info("order flow completed",
		zap.String("op", op),
		zap.Int64("telegram_id", req.TelegramID.Int64()),
		zap.Duration("duration", duration),
	)

	return nil
}
