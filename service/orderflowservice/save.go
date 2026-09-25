package orderflowservice

import (
	"context"
	"errors"
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

	if s.config.OrderTTL <= 0 {
		return richerror.New(
			op,
			errors.New("order flow TTL must be greater than zero"),
		).
			WithKind(richerror.KindInvalid).
			WithCode(richerror.CodeInvalidInput)
	}

	now := time.Now()

	if req.State.ExpiresAt <= 0 {
		req.State.ExpiresAt = now.Add(s.config.OrderTTL).Unix()
	}

	expiresAt := time.Unix(req.State.ExpiresAt, 0)
	remainingTTL := time.Until(expiresAt)

	if remainingTTL <= 0 {
		metrics.OrderFlowStateSaved.WithLabelValues(string(req.State.Stage), "expired").Inc()

		return richerror.New(op, errors.New("order flow has expired")).
			WithKind(richerror.KindConflict).
			WithCode(richerror.CodeInvalidInput)
	}

	req.TTLMins = remainingTTL

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
