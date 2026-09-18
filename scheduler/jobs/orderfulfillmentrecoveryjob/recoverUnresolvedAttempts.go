package orderfulfillmentrecoveryjob

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/richerror"
	"time"

	"go.uber.org/zap"
)

func (j *Job) recoverUnresolvedAttempts(ctx context.Context) error {
	const Op = "orderfulfillmentrecoveryjob.recoverUnresolvedAttempts"

	attempts, err := j.orderService.GetUnresolvedFulfillmentAttempts(
		ctx,
		j.config.BatchSize,
	)
	if err != nil {
		return richerror.New(Op, err)
	}

	if len(attempts) == 0 {
		return nil
	}

	for _, attempt := range attempts {
		switch attempt.Outcome {
		case orderentity.FulfillmentAttemptOutcomeCreated:
			if err := j.recoverCreatedAttempt(ctx, attempt); err != nil {
				metrics.WorkerRuns.
					WithLabelValues(j.Name(), "attempt_recovery_failed").
					Inc()

				logger.Logger.Error(
					"failed to recover created SMM fulfillment attempt",
					zap.Uint64("attempt_id", attempt.ID),
					zap.Uint64("order_id", attempt.OrderID),
					zap.Error(err),
				)
			}

		case orderentity.FulfillmentAttemptOutcomeUnknown:
			if err := j.recoverUnknownAttempt(ctx, attempt); err != nil {
				metrics.WorkerRuns.
					WithLabelValues(j.Name(), "unknown_attempt_recovery_failed").
					Inc()

				logger.Logger.Error(
					"failed to recover unknown SMM fulfillment attempt",
					zap.Uint64("attempt_id", attempt.ID),
					zap.Uint64("order_id", attempt.OrderID),
					zap.Error(err),
				)
			}
		}
	}

	return nil
}

func (j *Job) recoverCreatedAttempt(
	ctx context.Context,
	attempt orderentity.FulfillmentAttempt,
) error {
	const Op = "orderfulfillmentrecoveryjob.recoverCreatedAttempt"

	if attempt.ExternalOrderID == "" {
		return richerror.New(Op, nil).
			WithKind(richerror.KindInternal)
	}

	order, err := j.orderService.GetById(ctx, attempt.OrderID)
	if err != nil {
		return richerror.New(Op, err)
	}

	if order.ProviderID != nil &&
		*order.ProviderID != attempt.ProviderID {
		return richerror.New(Op, nil).
			WithKind(richerror.KindConflict)
	}

	if order.ExternalOrderID != "" {
		if order.ExternalOrderID != attempt.ExternalOrderID {
			return richerror.New(Op, nil).
				WithKind(richerror.KindConflict)
		}

		if err := j.orderService.MarkFulfillmentAttemptResolved(
			ctx,
			attempt.ID,
		); err != nil {
			return richerror.New(Op, err)
		}

		metrics.WorkerRuns.
			WithLabelValues(j.Name(), "attempt_already_persisted").
			Inc()

		return nil
	}

	if order.Status != orderentity.OrderStatusProcessing {
		logger.Logger.Error(
			"created fulfillment attempt belongs to non-processing order without external order id",
			zap.Uint64("attempt_id", attempt.ID),
			zap.Uint64("order_id", attempt.OrderID),
			zap.String("status", string(order.Status)),
		)

		return richerror.New(Op, nil).
			WithKind(richerror.KindConflict)
	}

	if err := j.orderService.SaveExternalOrder(
		ctx,
		attempt.OrderID,
		attempt.ProviderID,
		attempt.ExternalOrderID,
	); err != nil {
		logger.Logger.Error(
			"failed to persist recovered SMM provider order",
			zap.Uint64("attempt_id", attempt.ID),
			zap.Uint64("order_id", attempt.OrderID),
			zap.Uint64("provider_id", attempt.ProviderID),
			zap.String("external_order_id", attempt.ExternalOrderID),
			zap.Error(err),
		)

		metrics.WorkerRuns.
			WithLabelValues(j.Name(), "attempt_recovery_persist_failed").
			Inc()

		return richerror.New(Op, err)
	}

	if err := j.orderService.MarkFulfillmentAttemptResolved(
		ctx,
		attempt.ID,
	); err != nil {
		logger.Logger.Error(
			"failed to mark recovered fulfillment attempt resolved",
			zap.Uint64("attempt_id", attempt.ID),
			zap.Uint64("order_id", attempt.OrderID),
			zap.Error(err),
		)

		metrics.WorkerRuns.
			WithLabelValues(j.Name(), "attempt_resolve_failed").
			Inc()

		return richerror.New(Op, err)
	}

	logger.Logger.Info(
		"created SMM fulfillment attempt recovered",
		zap.Uint64("attempt_id", attempt.ID),
		zap.Uint64("order_id", attempt.OrderID),
		zap.Uint64("provider_id", attempt.ProviderID),
		zap.String("external_order_id", attempt.ExternalOrderID),
	)

	metrics.WorkerRuns.
		WithLabelValues(j.Name(), "attempt_recovered").
		Inc()

	return nil
}
func (j *Job) recoverUnknownAttempt(ctx context.Context, attempt orderentity.FulfillmentAttempt) error {
	const Op = "orderfulfillmentrecoveryjob.recoverUnknownAttempt"

	order, err := j.orderService.GetById(ctx, attempt.OrderID)
	if err != nil {
		return richerror.New(Op, err)
	}

	if order.ProviderID == nil {
		if err := j.orderService.AssignProvider(
			ctx,
			attempt.OrderID,
			attempt.ProviderID,
		); err != nil {
			return richerror.New(Op, err)
		}

		logger.Logger.Info(
			"provider assignment recovered for unknown SMM attempt",
			zap.Uint64("attempt_id", attempt.ID),
			zap.Uint64("order_id", attempt.OrderID),
			zap.Uint64("provider_id", attempt.ProviderID),
		)

		metrics.WorkerRuns.WithLabelValues(j.Name(), "unknown_attempt_provider_assigned").Inc()

		return nil
	}

	if *order.ProviderID != attempt.ProviderID {
		return richerror.New(Op, nil).
			WithKind(richerror.KindConflict)
	}

	age := time.Since(attempt.CreatedAt)

	if age < 0 {
		logger.Logger.Warn(
			"unknown SMM fulfillment attempt has future creation time",
			zap.Uint64("attempt_id", attempt.ID),
			zap.Uint64("order_id", attempt.OrderID),
			zap.Uint64("provider_id", attempt.ProviderID),
			zap.Time("created_at", attempt.CreatedAt),
		)

		metrics.WorkerRuns.WithLabelValues(j.Name(), "unknown_attempt_pending").Inc()

		return nil
	}

	if j.config.ReconciliationAfter > 0 && age >= j.config.ReconciliationAfter {
		metrics.WorkerRuns.WithLabelValues(j.Name(), "reconciliation_required").Inc()

		logger.Logger.Error(
			"unknown SMM fulfillment attempt requires reconciliation",
			zap.Uint64("attempt_id", attempt.ID),
			zap.Uint64("order_id", attempt.OrderID),
			zap.Uint64("provider_id", attempt.ProviderID),
			zap.Duration("age", age),
			zap.Time("created_at", attempt.CreatedAt),
		)

		return nil
	}

	metrics.WorkerRuns.WithLabelValues(j.Name(), "unknown_attempt_pending").Inc()

	return nil
}
