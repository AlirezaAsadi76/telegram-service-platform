package checkoutservice

import (
	"context"
	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/params/notificationparams"
	"time"

	"telegram-service-platform/logger"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/walletparam"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"go.uber.org/zap"
)

func (s *Service) ProcessWalletPurchase(ctx context.Context, req checkoutparams.WalletPurchaseRequest) error {
	const Op = "checkoutservice.ProcessWalletPurchase"

	start := time.Now()

	logger.Logger.Info(
		"wallet purchase started",
		zap.Uint64("user_id", req.UserID),
		zap.String("idempotency_key", req.IdempotencyKey),
		zap.String("amount", req.Amount.String()),
	)

	result, err := s.walletPurchaseRepo.ExecuteWalletPurchase(
		ctx,
		walletparam.WalletPurchaseRequest{
			UserID:         req.UserID,
			ProductType:    req.ProductType,
			ProductID:      req.ProductID,
			Quantity:       req.Quantity,
			TargetLink:     req.TargetLink,
			Amount:         req.Amount,
			Currency:       req.Currency,
			IdempotencyKey: req.IdempotencyKey,
		},
	)
	if err != nil {
		metrics.WalletTransactions.
			WithLabelValues("purchase_failed").
			Inc()

		metrics.CheckoutLatency.
			WithLabelValues("wallet").
			Observe(time.Since(start).Seconds())

		logger.Logger.Error(
			"wallet purchase transaction failed",
			zap.Uint64("user_id", req.UserID),
			zap.String("idempotency_key", req.IdempotencyKey),
			zap.Error(err),
			zap.Duration("latency", time.Since(start)),
		)

		return richerror.New(Op, err)
	}

	if s.fulfillmentEnqueuer == nil {
		logger.Logger.Error(
			"wallet purchase committed but fulfillment enqueuer is unavailable",
			zap.Uint64("order_id", result.OrderID),
		)

		return richerror.New(Op, nil).
			WithKind(richerror.KindInternal).
			WithMessage(msgerror.OrderFulfillmentEnqueueFailed)
	}

	if err := s.fulfillmentEnqueuer.Enqueue(ctx, result.OrderID); err != nil {
		metrics.WalletTransactions.
			WithLabelValues("enqueue_failed").
			Inc()

		logger.Logger.Error(
			"wallet purchase committed but fulfillment enqueue failed",
			zap.Uint64("order_id", result.OrderID),
			zap.Uint64("wallet_transaction_id", result.WalletTxID),
			zap.Error(err),
		)

		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.OrderFulfillmentEnqueueFailed)
	}

	metrics.WalletTransactions.
		WithLabelValues("purchase_success").
		Inc()

	metrics.OrdersTotal.
		WithLabelValues("wallet", "paid").
		Inc()

	metrics.CheckoutLatency.
		WithLabelValues("wallet").
		Observe(time.Since(start).Seconds())

	logger.Logger.Info(
		"wallet purchase completed",
		zap.Uint64("order_id", result.OrderID),
		zap.Uint64("wallet_transaction_id", result.WalletTxID),
		zap.Duration("latency", time.Since(start)),
	)

	if err := s.notificationSvc.Create(
		ctx,
		notificationparams.CreateRequest{
			UserID: req.UserID,
			Type:   notificationentity.NotificationTypeOrderPaid,
			Payload: map[string]any{
				"order_id": result.OrderID,
			},
		},
	); err != nil {
		logger.Logger.Error(
			"wallet purchase notification creation failed",
			zap.Uint64("order_id", result.OrderID),
			zap.Uint64("user_id", req.UserID),
			zap.Error(err),
		)
	}

	return nil
}
