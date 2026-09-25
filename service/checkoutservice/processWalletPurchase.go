package checkoutservice

import (
	"context"
	"errors"
	"telegram-service-platform/pkg/msgerror"
	"time"

	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/params/walletparam"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/richerror"

	"go.uber.org/zap"
)

func (s *Service) ProcessWalletPurchase(ctx context.Context, req checkoutparams.WalletPurchaseRequest) (*checkoutparams.WalletPurchaseResponse, error) {
	const Op = "checkoutservice.ProcessWalletPurchase"

	start := time.Now()

	if s.transactionRepo == nil {
		return nil, richerror.New(Op, errors.New("wallet transaction repository is unavailable")).
			WithKind(richerror.KindDependencyFailure).
			WithMessage(msgerror.ExternalServiceFailed)
	}

	logger.Logger.Info(
		"wallet purchase started",
		zap.Uint64("user_id", req.UserID),
		zap.String("idempotency_key", req.IdempotencyKey),
		zap.String("amount", req.Amount.String()),
	)

	result, err := s.transactionRepo.ExecuteWalletPurchase(
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

		return nil, richerror.New(Op, err)
	}

	if result == nil {
		return nil, richerror.New(Op, errors.New("wallet purchase returned nil result")).
			WithKind(richerror.KindInternal).
			WithCode(richerror.CodeUnknown)
	}

	logger.Logger.Info(
		"wallet purchase committed",
		zap.Uint64("order_id", result.OrderID),
		zap.Uint64("wallet_transaction_id", result.WalletTxID),
	)

	if s.notificationSvc != nil {
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
				zap.Error(err),
			)
		}
	}

	if s.fulfillmentEnqueuer == nil {
		metrics.WalletTransactions.
			WithLabelValues("enqueue_failed").
			Inc()

		logger.Logger.Error(
			"wallet purchase committed but fulfillment enqueuer is unavailable",
			zap.Uint64("order_id", result.OrderID),
			zap.Uint64("wallet_transaction_id", result.WalletTxID),
		)
	} else if err := s.fulfillmentEnqueuer.Enqueue(ctx, result.OrderID); err != nil {
		metrics.WalletTransactions.
			WithLabelValues("enqueue_failed").
			Inc()

		logger.Logger.Error(
			"wallet purchase committed but fulfillment enqueue failed; recovery will retry",
			zap.Uint64("order_id", result.OrderID),
			zap.Uint64("wallet_transaction_id", result.WalletTxID),
			zap.Error(err),
		)
	} else {
		metrics.WalletTransactions.
			WithLabelValues("enqueue_success").
			Inc()

		logger.Logger.Info(
			"wallet fulfillment enqueued",
			zap.Uint64("order_id", result.OrderID),
			zap.Uint64("wallet_transaction_id", result.WalletTxID),
		)
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

	return &checkoutparams.WalletPurchaseResponse{
		OrderID:    result.OrderID,
		WalletTxID: result.WalletTxID,
		Amount:     req.Amount,
		Currency:   req.Currency,
		NewBalance: result.NewBalance,
	}, nil
}
