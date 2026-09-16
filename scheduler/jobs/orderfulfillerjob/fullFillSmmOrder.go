package orderfulfillerjob

import (
	"context"
	"fmt"
	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/params/smmparams"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"go.uber.org/zap"
)

func (j *Job) fulfillSMMOrder(ctx context.Context, order *orderentity.Order) error {
	const Op = "orderfulfillerjob.fulfillSMMOrder"

	result, err := j.smmProviderService.CreateOrder(
		ctx,
		smmparams.CreateOrderAdapterRequest{
			ServiceID: fmt.Sprintf("%d", order.ProductID),
			Link:      order.TargetLink,
			Quantity:  order.Quantity,
		},
	)

	if err != nil {
		if result.ProviderID != 0 {
			logger.Logger.Error(
				"SMM create returned unknown result",
				zap.Uint64("order_id", order.ID),
				zap.Uint64("provider_id", result.ProviderID),
				zap.String("provider_name", result.ProviderName),
				zap.String("outcome", string(result.Outcome)),
				zap.Error(err),
			)
		}

		metrics.SMMProviderRequests.WithLabelValues(result.ProviderName, "error").Inc()

		if result.Outcome == smmparams.CreateOrderOutcomeUnknown {
			logger.Logger.Warn(
				"SMM create outcome is unknown; order remains processing",
				zap.Uint64("order_id", order.ID),
				zap.Uint64("provider_id", result.ProviderID),
				zap.String("provider_name", result.ProviderName),
			)

			metrics.WorkerRuns.
				WithLabelValues(j.Name(), "unknown_provider_result").
				Inc()

			return nil
		}

		return richerror.New(Op, err)
	}

	switch result.Outcome {
	case smmparams.CreateOrderOutcomeCreated:
		if result.ProviderID == 0 ||
			result.ExternalOrderID == "" {
			logger.Logger.Error(
				"SMM provider returned incomplete created result",
				zap.Uint64("order_id", order.ID),
				zap.Uint64("provider_id", result.ProviderID),
				zap.String("external_order_id", result.ExternalOrderID),
			)

			return richerror.New(Op, nil).
				WithKind(richerror.KindExternalAPI).
				WithCode(richerror.CodeSMMProviderInvalidResponse).
				WithMessage(msgerror.SMMProviderInvalidResponse)
		}

		if err := j.orderService.SetProviderOrder(
			ctx,
			order.ID,
			result.ProviderID,
			result.ExternalOrderID,
		); err != nil {
			logger.Logger.Error(
				"failed to persist SMM provider order",
				zap.Uint64("order_id", order.ID),
				zap.Uint64("provider_id", result.ProviderID),
				zap.String("external_order_id", result.ExternalOrderID),
				zap.Error(err),
			)

			metrics.WorkerRuns.
				WithLabelValues(j.Name(), "provider_result_persist_failed").
				Inc()

			// Never retry Create blindly after the provider has
			// already returned an external order id.
			return nil
		}

		metrics.SMMProviderRequests.
			WithLabelValues(result.ProviderName, "success").
			Inc()

		logger.Logger.Info(
			"SMM order created successfully",
			zap.Uint64("order_id", order.ID),
			zap.Uint64("provider_id", result.ProviderID),
			zap.String("provider_name", result.ProviderName),
			zap.String("external_order_id", result.ExternalOrderID),
		)

		if err := j.notificationService.Create(
			ctx,
			notificationparams.CreateRequest{
				UserID: order.UserID,
				Type:   notificationentity.NotificationTypeOrderProcessing,
				Payload: map[string]any{
					"order_id": order.ID,
				},
			},
		); err != nil {
			logger.Logger.Error(
				"failed to create order processing notification",
				zap.Uint64("order_id", order.ID),
				zap.Error(err),
			)
		}

		metrics.WorkerRuns.
			WithLabelValues(j.Name(), "success").
			Inc()

		return nil

	case smmparams.CreateOrderOutcomeRejected:
		logger.Logger.Warn(
			"SMM providers rejected order",
			zap.Uint64("order_id", order.ID),
		)

		metrics.WorkerRuns.
			WithLabelValues(j.Name(), "provider_rejected").
			Inc()

		return nil

	case smmparams.CreateOrderOutcomeUnknown:
		logger.Logger.Warn(
			"SMM create result is unknown",
			zap.Uint64("order_id", order.ID),
			zap.Uint64("provider_id", result.ProviderID),
			zap.String("provider_name", result.ProviderName),
		)

		metrics.WorkerRuns.
			WithLabelValues(j.Name(), "unknown_provider_result").
			Inc()

		return nil

	default:
		return richerror.New(Op, fmt.Errorf("unsupported fulfillment outcome: %s", result.Outcome)).
			WithKind(richerror.KindInternal).
			WithMessage(msgerror.InternalServerError)
	}
}
