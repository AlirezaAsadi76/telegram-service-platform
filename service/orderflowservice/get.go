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

	if state == nil {
		metrics.OrderFlowStateRetrieved.WithLabelValues("false").Inc()

		logger.Logger.Debug(
			"order flow state not found",
			zap.String("op", op),
			zap.Int64("telegram_id", req.TelegramID.Int64()),
			zap.Duration("latency", time.Since(start)),
		)

		return nil, nil
	}

	if state.ExpiresAt > 0 && time.Now().Unix() >= state.ExpiresAt {
		metrics.OrderFlowStateRetrieved.WithLabelValues("false").Inc()

		deleteErr := s.DeleteOrderFlow(
			ctx,
			orderparams.DeleteOrderFlowRequest{
				TelegramID: req.TelegramID,
			},
			"expired",
		)

		if deleteErr != nil {
			logger.Logger.Error(
				"failed to delete expired order flow state",
				zap.String("op", op),
				zap.Int64("telegram_id", req.TelegramID.Int64()),
				zap.Int64("expires_at", state.ExpiresAt),
				zap.Error(deleteErr),
			)
		}

		logger.Logger.Info(
			"expired order flow state rejected",
			zap.String("op", op),
			zap.Int64("telegram_id", req.TelegramID.Int64()),
			zap.Int64("expires_at", state.ExpiresAt),
			zap.Duration("latency", time.Since(start)),
		)

		return nil, nil
	}

	metrics.OrderFlowStateRetrieved.WithLabelValues("true").Inc()

	logger.Logger.Debug(
		"order flow state retrieved",
		zap.String("op", op),
		zap.Int64("telegram_id", req.TelegramID.Int64()),
		zap.String("stage", string(state.Stage)),
		zap.Int64("expires_at", state.ExpiresAt),
		zap.Duration("latency", time.Since(start)),
	)

	return state, nil
}
