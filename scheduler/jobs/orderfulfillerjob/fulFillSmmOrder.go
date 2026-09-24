package orderfulfillerjob

import (
	"context"
	"fmt"
	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/params/smmparams"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"go.uber.org/zap"
)

const smmProviderRejectedRefundReason = "smm_provider_rejected"

func (j *Job) fulfillSMMOrder(ctx context.Context, order *orderentity.Order) error {
	const Op = "orderfulfillerjob.fulfillSMMOrder"

	if j.productService == nil {
		return richerror.New(Op, fmt.Errorf("product service is not configured")).
			WithKind(richerror.KindDependencyFailure).
			WithMessage(msgerror.ExternalServiceFailed)
	}

	smm, err := j.productService.GetSMMServiceByMappingID(ctx, int64(order.ProductID))
	if err != nil {
		metrics.WorkerRuns.WithLabelValues(j.Name(), "smm_service_resolution_failed").Inc()

		logger.Logger.Error(
			"failed to resolve SMM service from order mapping",
			zap.Uint64("order_id", order.ID),
			zap.Uint64("mapping_id", order.ProductID),
			zap.Error(err),
		)

		return richerror.New(Op, err).
			WithKind(richerror.KindDependencyFailure).
			WithMessage(msgerror.ExternalServiceFailed)
	}

	if smm == nil {
		metrics.WorkerRuns.WithLabelValues(j.Name(), "smm_service_resolution_failed").Inc()

		return richerror.New(Op, fmt.Errorf("resolved SMM service is nil")).
			WithKind(richerror.KindDependencyFailure).
			WithMessage(msgerror.ExternalServiceFailed)
	}

	if smm.Service <= 0 {
		metrics.WorkerRuns.WithLabelValues(j.Name(), "invalid_smm_service").Inc()

		logger.Logger.Error(
			"resolved SMM service has invalid provider service ID",
			zap.Uint64("order_id", order.ID),
			zap.Int64("smm_id", smm.Id),
			zap.Int64("provider_service_id", smm.Service),
			zap.String("provider_name", smm.ProviderName),
		)

		return richerror.New(Op, fmt.Errorf("invalid provider service ID for SMM %d", smm.Id)).
			WithKind(richerror.KindInvalid).
			WithCode(richerror.CodeSMMProviderInvalidResponse).
			WithMessage(msgerror.SMMProviderInvalidResponse)
	}

	if smm.ProviderName == "" {
		metrics.WorkerRuns.WithLabelValues(j.Name(), "invalid_smm_service").Inc()

		logger.Logger.Error(
			"resolved SMM service has no provider name",
			zap.Uint64("order_id", order.ID),
			zap.Int64("smm_id", smm.Id),
			zap.Int64("provider_service_id", smm.Service),
		)

		return richerror.New(Op, fmt.Errorf("provider name is missing for SMM %d", smm.Id)).
			WithKind(richerror.KindInvalid).
			WithCode(richerror.CodeSMMProviderInvalidResponse).
			WithMessage(msgerror.SMMProviderInvalidResponse)
	}

	logger.Logger.Debug(
		"resolved SMM fulfillment identity",
		zap.Uint64("order_id", order.ID),
		zap.Uint64("mapping_id", order.ProductID),
		zap.Int64("smm_id", smm.Id),
		zap.Int64("provider_service_id", smm.Service),
		zap.String("provider_name", smm.ProviderName),
	)

	result, scErr := j.smmProviderService.CreateOrder(
		ctx,
		smmparams.CreateOrderAdapterRequest{
			ProviderName: smm.ProviderName,
			ServiceID:    fmt.Sprintf("%d", smm.Service),
			Link:         order.TargetLink,
			Quantity:     order.Quantity,
		},
	)

	if scErr != nil {
		if result != nil {
			metrics.SMMProviderRequests.WithLabelValues(smm.ProviderName, "error").Inc()

			if result.Outcome == smmparams.CreateOrderOutcomeUnknown {
				logger.Logger.Warn(
					"SMM create outcome is unknown",
					zap.Uint64("order_id", order.ID),
					zap.Uint64("provider_id", result.ProviderID),
					zap.String("provider_name", result.ProviderName),
					zap.Error(scErr),
				)

				if result.ProviderID != 0 {
					attempt := &orderentity.FulfillmentAttempt{
						OrderID:    order.ID,
						ProviderID: result.ProviderID,
						Outcome:    orderentity.FulfillmentAttemptOutcomeUnknown,
					}

					if _, attemptErr := j.orderService.CreateFulfillmentAttempt(ctx, attempt); attemptErr != nil {
						logger.Logger.Error(
							"failed to persist unknown SMM fulfillment attempt",
							zap.Uint64("order_id", order.ID),
							zap.Uint64("provider_id", result.ProviderID),
							zap.Error(attemptErr),
						)

						metrics.WorkerRuns.
							WithLabelValues(j.Name(), "fulfillment_attempt_persist_failed").
							Inc()

						return richerror.New(Op, attemptErr).
							WithKind(richerror.KindQueryFailure).
							WithMessage(msgerror.OrderUpdateFailed)
					}
					if providerErr := j.orderService.AssignProvider(ctx, order.ID, result.ProviderID); providerErr != nil {
						logger.Logger.Error(
							"failed to persist provider for unknown SMM result",
							zap.Uint64("order_id", order.ID),
							zap.Uint64("provider_id", result.ProviderID),
							zap.Error(providerErr),
						)

						metrics.WorkerRuns.WithLabelValues(j.Name(), "provider_persist_failed").Inc()

						return richerror.New(Op, providerErr).
							WithKind(richerror.KindQueryFailure).
							WithMessage(msgerror.OrderUpdateFailed)
					}
				}

				metrics.WorkerRuns.WithLabelValues(j.Name(), "unknown_provider_result").Inc()

				return nil
			}
		}

		return richerror.New(Op, scErr)
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

		if err := j.orderService.SaveExternalOrder(
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

			// Provider has already created the order.
			// Never retry Create blindly.
			return nil
		}

		attempt := &orderentity.FulfillmentAttempt{
			OrderID:         order.ID,
			ProviderID:      result.ProviderID,
			Outcome:         orderentity.FulfillmentAttemptOutcomeCreated,
			ExternalOrderID: result.ExternalOrderID,
		}

		attemptID, attemptErr := j.orderService.CreateFulfillmentAttempt(
			ctx,
			attempt,
		)
		if attemptErr != nil {
			logger.Logger.Error(
				"failed to persist created SMM fulfillment attempt after provider result was saved",
				zap.Uint64("order_id", order.ID),
				zap.Uint64("provider_id", result.ProviderID),
				zap.String("external_order_id", result.ExternalOrderID),
				zap.Error(attemptErr),
			)

			metrics.WorkerRuns.
				WithLabelValues(j.Name(), "fulfillment_attempt_persist_failed").
				Inc()

			// Order provider data is already durable.
			// StatusSync can continue reconciliation.
			return nil
		}

		attempt.ID = attemptID

		if err := j.orderService.MarkFulfillmentAttemptResolved(
			ctx,
			attempt.ID,
		); err != nil {
			logger.Logger.Error(
				"failed to mark SMM fulfillment attempt resolved",
				zap.Uint64("order_id", order.ID),
				zap.Uint64("attempt_id", attempt.ID),
				zap.Error(err),
			)

			metrics.WorkerRuns.
				WithLabelValues(j.Name(), "fulfillment_attempt_resolve_failed").
				Inc()
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

		if err := j.checkoutService.RefundOrder(
			ctx,
			checkoutparams.RefundOrderRequest{
				OrderID: order.ID,
				Reason:  smmProviderRejectedRefundReason,
			},
		); err != nil {
			metrics.WorkerRuns.
				WithLabelValues(j.Name(), "refund_failed").
				Inc()

			logger.Logger.Error(
				"failed to refund rejected SMM order",
				zap.Uint64("order_id", order.ID),
				zap.Error(err),
			)

			return richerror.New(Op, err)
		}

		metrics.WorkerRuns.
			WithLabelValues(j.Name(), "provider_rejected_refunded").
			Inc()

		logger.Logger.Info(
			"rejected SMM order refunded successfully",
			zap.Uint64("order_id", order.ID),
		)

		return nil

	case smmparams.CreateOrderOutcomeUnknown:
		logger.Logger.Warn(
			"SMM create result is unknown",
			zap.Uint64("order_id", order.ID),
			zap.Uint64("provider_id", result.ProviderID),
			zap.String("provider_name", result.ProviderName),
		)

		if result.ProviderID != 0 {
			attempt := &orderentity.FulfillmentAttempt{
				OrderID:    order.ID,
				ProviderID: result.ProviderID,
				Outcome:    orderentity.FulfillmentAttemptOutcomeUnknown,
			}

			if _, attemptErr := j.orderService.CreateFulfillmentAttempt(ctx, attempt); attemptErr != nil {
				logger.Logger.Error(
					"failed to persist unknown SMM fulfillment attempt",
					zap.Uint64("order_id", order.ID),
					zap.Uint64("provider_id", result.ProviderID),
					zap.Error(attemptErr),
				)

				metrics.WorkerRuns.
					WithLabelValues(j.Name(), "fulfillment_attempt_persist_failed").
					Inc()

				return richerror.New(Op, attemptErr).
					WithKind(richerror.KindQueryFailure).
					WithMessage(msgerror.OrderUpdateFailed)
			}

			if providerErr := j.orderService.AssignProvider(ctx, order.ID, result.ProviderID); providerErr != nil {
				logger.Logger.Error(
					"failed to persist provider for unknown SMM result",
					zap.Uint64("order_id", order.ID),
					zap.Uint64("provider_id", result.ProviderID),
					zap.Error(providerErr),
				)

				return richerror.New(Op, providerErr).
					WithKind(richerror.KindQueryFailure).
					WithMessage(msgerror.OrderUpdateFailed)
			}
		}

		metrics.WorkerRuns.WithLabelValues(j.Name(), "unknown_provider_result").Inc()

		return nil

	default:
		return richerror.New(Op, fmt.Errorf("unsupported fulfillment outcome: %s", result.Outcome)).
			WithKind(richerror.KindInternal).
			WithMessage(msgerror.InternalServerError)
	}
}
