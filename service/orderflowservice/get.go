package orderflowservice

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/params/orderparams"
	"time"

	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/richerror"

	"telegram-service-platform/logger"

	"go.uber.org/zap"
)

func (s *Service) GetOrderFlow(ctx context.Context, req orderparams.GetOrderFlowRequest) (*orderentity.OrderFlowState, error) {
	const op = "orderflowservice.GetOrderFlow"
	start := time.Now()

	if req.TelegramID <= 0 {
		return nil, richerror.New(op, nil).
			WithKind(richerror.KindValidation)
	}

	state, err := s.repo.Get(ctx, req)
	if err != nil {
		metrics.OrderFlowStateRetrieved.WithLabelValues("error").Inc()
		logger.Logger.Error("failed to get order flow state",
			zap.String("op", op),
			zap.Int64("telegram_id", req.TelegramID.Int64()),
			zap.Error(err),
		)
		return nil, richerror.New(op, err)
	}

	found := state != nil
	metrics.OrderFlowStateRetrieved.WithLabelValues(boolToStr(found)).Inc()

	if found {
		logger.Logger.Debug("order flow state retrieved",
			zap.String("op", op),
			zap.Int64("telegram_id", req.TelegramID.Int64()),
			zap.String("stage", string(state.Stage)),
			zap.Duration("latency", time.Since(start)),
		)
	} else {
		logger.Logger.Debug("order flow state not found",
			zap.String("op", op),
			zap.Int64("telegram_id", req.TelegramID.Int64()),
			zap.Duration("latency", time.Since(start)),
		)
	}

	return state, nil
}

func boolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
