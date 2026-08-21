package orderflowservice

import (
	"context"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/richerror"
	"time"

	"go.uber.org/zap"
)

func (s *Service) SaveOrderFlow(ctx context.Context, req orderparams.SaveOrderFlowRequest) error {
	const op = "orderflowservice.SaveOrderFlow"
	start := time.Now()

	if req.TelegramID <= 0 {
		return richerror.New(op, nil).
			WithKind(richerror.KindValidation)
	}

	if err := s.repo.Save(ctx, req); err != nil {
		metrics.OrderFlowStateSaved.WithLabelValues(string(req.State.Stage), "error").Inc()
		logger.Logger.Error("failed to save order flow state",
			zap.String("op", op),
			zap.Int64("telegram_id", req.TelegramID.Int64()),
			zap.Uint64("service_id", req.State.ServiceID),
			zap.String("stage", string(req.State.Stage)),
			zap.Error(err),
		)
		return richerror.New(op, err)
	}

	metrics.OrderFlowStateSaved.WithLabelValues(string(req.State.Stage), "success").Inc()
	logger.Logger.Info("order flow state saved",
		zap.String("op", op),
		zap.Int64("telegram_id", req.TelegramID.Int64()),
		zap.Uint64("service_id", req.State.ServiceID),
		zap.String("stage", string(req.State.Stage)),
		zap.Duration("latency", time.Since(start)),
	)

	return nil
}
