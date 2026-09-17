package checkoutservice

import (
	"context"
	"time"

	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/richerror"

	"telegram-service-platform/logger"

	"go.uber.org/zap"
)

func (s *Service) RefundOrder(ctx context.Context, req checkoutparams.RefundOrderRequest) error {
	const Op = "checkoutservice.RefundOrder"

	start := time.Now()

	logger.Logger.Info(
		"order refund started",
		zap.Uint64("order_id", req.OrderID),
		zap.String("reason", req.Reason),
	)

	result, err := s.transactionRepo.ExecuteOrderRefund(ctx, req)
	if err != nil {
		metrics.WalletTransactions.WithLabelValues("refund_failed").Inc()

		logger.Logger.Error(
			"order refund failed",
			zap.Uint64("order_id", req.OrderID),
			zap.Error(err),
		)

		return richerror.New(Op, err)
	}

	metrics.WalletTransactions.
		WithLabelValues("refund_success").
		Inc()

	logger.Logger.Info(
		"order refund completed",
		zap.Uint64("order_id", result.OrderID),
		zap.Uint64("wallet_transaction_id", result.WalletTxID),
		zap.String("amount", result.RefundAmount.String()),
		zap.Bool("already_refunded", result.AlreadyRefunded),
		zap.Duration("latency", time.Since(start)),
	)

	if err := s.notificationSvc.Create(
		ctx,
		notificationparams.CreateRequest{
			UserID: result.UserID,
			Type:   notificationentity.NotificationTypeOrderFailed,
			Payload: map[string]any{
				"order_id": result.OrderID,
				"reason":   req.Reason,
			},
		},
	); err != nil {
		logger.Logger.Error(
			"refund notification failed",
			zap.Uint64("order_id", result.OrderID),
			zap.Error(err),
		)
	}

	return nil
}
