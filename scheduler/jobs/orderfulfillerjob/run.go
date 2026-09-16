package orderfulfillerjob

import (
	"context"
	"errors"
	"fmt"
	"telegram-service-platform/params/orderparams"
	"time"

	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/params/smmparams"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/pkg/unmarshal"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (j *Job) Run(ctx context.Context) error {
	const Op = "orderfulfillerjob.Run"

	start := time.Now()
	jobName := j.Name()

	defer func() {
		metrics.WorkerDuration.
			WithLabelValues(jobName).
			Observe(time.Since(start).Seconds())
	}()

	logger.Logger.Info(
		"order fulfillment worker started",
		zap.String("job", jobName),
	)

	j.mutex.Lock()
	result, err := j.redis.BRPop(ctx, j.config.Timeout, j.config.QueueKey)
	j.mutex.Unlock()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			logger.Logger.Debug(
				"order fulfillment queue is empty",
				zap.String("job", jobName),
			)
			return nil
		}

		metrics.WorkerRuns.WithLabelValues(jobName, "error").Inc()

		logger.Logger.Error(
			"order fulfillment worker redis error",
			zap.String("job", jobName),
			zap.Error(err),
		)

		return richerror.New(Op, err).
			WithKind(richerror.KindDependencyFailure).
			WithMessage(msgerror.ExternalServiceFailed)
	}

	if len(result) < 2 {
		metrics.WorkerRuns.
			WithLabelValues(jobName, "invalid_message").
			Inc()

		logger.Logger.Error(
			"invalid order fulfillment queue message",
			zap.String("job", jobName),
			zap.Int("items", len(result)),
		)

		return richerror.New(Op, fmt.Errorf("invalid queue message")).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	orderID, uErr := unmarshal.UnmarshalToUint64(result[1])
	if uErr != nil {
		metrics.WorkerRuns.
			WithLabelValues(jobName, "invalid_message").
			Inc()

		logger.Logger.Error(
			"failed to unmarshal order id",
			zap.String("job", jobName),
			zap.String("payload", result[1]),
			zap.Error(err),
		)

		return richerror.New(Op, uErr).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.UnmarshalFailed)
	}

	order, gErr := j.orderService.GetById(ctx, orderID)
	if gErr != nil {
		logger.Logger.Error(
			"failed to get order",
			zap.String("job", jobName),
			zap.Uint64("order_id", orderID),
			zap.Error(gErr),
		)

		metrics.WorkerRuns.
			WithLabelValues(jobName, "order_not_found").
			Inc()

		return nil
	}

	if order.Status != orderentity.OrderStatusPaid {
		logger.Logger.Info(
			"order is no longer paid; skipping queue message",
			zap.String("job", jobName),
			zap.Uint64("order_id", order.ID),
			zap.String("status", string(order.Status)),
		)

		metrics.WorkerRuns.WithLabelValues(jobName, "already_processed").Inc()

		return nil
	}

	claimed, cErr := j.orderService.ClaimForProcessing(ctx, order.ID)
	if cErr != nil {
		metrics.WorkerRuns.
			WithLabelValues(jobName, "claim_failed").
			Inc()

		logger.Logger.Error(
			"failed to claim order",
			zap.String("job", jobName),
			zap.Uint64("order_id", order.ID),
			zap.Error(err),
		)

		return cErr
	}

	if !claimed {
		logger.Logger.Info(
			"order was already claimed",
			zap.String("job", jobName),
			zap.Uint64("order_id", order.ID),
		)

		metrics.WorkerRuns.
			WithLabelValues(jobName, "already_claimed").
			Inc()

		return nil
	}

	logger.Logger.Info(
		"order claimed for fulfillment",
		zap.String("job", jobName),
		zap.Uint64("order_id", order.ID),
	)

	switch order.ProductType {
	case productentity.ProductTypeSMM:
		return j.fulfillSMMOrder(ctx, order)

	default:
		logger.Logger.Error(
			"unsupported order product type",
			zap.String("job", jobName),
			zap.Uint64("order_id", order.ID),
			zap.String("product_type", string(order.ProductType)),
		)

		metrics.WorkerRuns.
			WithLabelValues(jobName, "unsupported_product_type").
			Inc()

		return j.markUnsupportedOrder(
			ctx,
			order.ID,
		)
	}
}

func (j *Job) markUnsupportedOrder(ctx context.Context, orderID uint64) error {
	const Op = "orderfulfillerjob.markUnsupportedOrder"

	status := orderentity.OrderStatusFailed

	if err := j.orderService.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
		OrderID: orderID,
		Status:  status,
	}); err != nil {
		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.OrderUpdateFailed)
	}

	return nil
}
