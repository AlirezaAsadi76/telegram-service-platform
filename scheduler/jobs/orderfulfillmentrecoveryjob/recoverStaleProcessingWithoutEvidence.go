package orderfulfillmentrecoveryjob

import (
	"context"
	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/richerror"
	"time"

	"go.uber.org/zap"
)

func (j *Job) recoverStaleProcessingWithoutEvidence(ctx context.Context) error {
	const Op = "orderfulfillmentrecoveryjob.recoverStaleProcessingWithoutEvidence"

	if j.config.ReconciliationAfter <= 0 {
		return nil
	}

	orders, err := j.orderService.GetStaleProcessingWithoutEvidence(
		ctx,
		j.config.ReconciliationAfter,
		j.config.BatchSize,
	)
	if err != nil {
		return richerror.New(Op, err)
	}

	for _, order := range orders {
		metrics.WorkerRuns.
			WithLabelValues(j.Name(), "reconciliation_required").
			Inc()

		logger.Logger.Error(
			"processing SMM order has no durable fulfillment evidence and requires reconciliation",
			zap.Uint64("order_id", order.ID),
			zap.Uint64("user_id", order.UserID),
			zap.String("status", string(order.Status)),
			zap.Time("updated_at", order.UpdatedAt),
			zap.Duration("age", time.Since(order.UpdatedAt)),
		)
	}

	return nil
}
