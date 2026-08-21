package orderflowservice

import (
	"context"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/metrics"
	"time"

	"go.uber.org/zap"
)

func (s *Service) AbandonOrderFlow(ctx context.Context, req orderparams.DeleteOrderFlowRequest, duration time.Duration) error {
	const op = "orderflowservice.AbandonOrderFlow"

	// ثبت متریک انصراف
	metrics.OrderFlowDuration.WithLabelValues("abandoned").Observe(duration.Seconds())

	// حذف state
	if err := s.DeleteOrderFlow(ctx, req, "abandoned"); err != nil {
		logger.Logger.Error("failed to abandon order flow",
			zap.String("op", op),
			zap.Int64("telegram_id", req.TelegramID.Int64()),
			zap.Error(err),
		)
		return err
	}

	logger.Logger.Info("order flow abandoned",
		zap.String("op", op),
		zap.Int64("telegram_id", req.TelegramID.Int64()),
		zap.Duration("duration", duration),
	)

	return nil
}
