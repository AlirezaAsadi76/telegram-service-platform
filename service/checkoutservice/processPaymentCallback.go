package checkoutservice

import (
	"context"
	"fmt"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"

	"go.uber.org/zap"
)

func (s *Service) ProcessPaymentCallback(ctx context.Context, paymentID uint64, externalID string, callbackData map[string]any) error {
	const Op = "checkoutservice.ProcessPaymentCallback"

	start := time.Now()
	logger.Logger.Info("payment callback started",
		zap.Uint64("payment_id", paymentID),
	)

	// 1. Verify payment
	verifyResp, vErr := s.paymentSvc.Verify(ctx, paymentparams.VerifyRequest{
		PaymentID:    paymentID,
		ExternalID:   externalID,
		CallbackData: callbackData,
	})
	if vErr != nil {
		metrics.PaymentsProcessed.WithLabelValues("gateway", "failed").Inc()
		metrics.CheckoutLatency.WithLabelValues("payment_callback").Observe(time.Since(start).Seconds())
		logger.Logger.Error("payment callback failed", zap.Error(vErr), zap.Uint64("payment_id", paymentID))
		return richerror.New(Op, vErr)
	}

	// Get payment details
	payment, err := s.paymentSvc.GetById(ctx, paymentID)
	if err != nil {
		metrics.PaymentsProcessed.WithLabelValues("gateway", "failed").Inc()
		metrics.CheckoutLatency.WithLabelValues("payment_callback").Observe(time.Since(start).Seconds())
		logger.Logger.Error("payment callback failed", zap.Error(vErr), zap.Uint64("payment_id", paymentID))
		return richerror.New(Op, err)
	}

	if verifyResp.Status == paymententity.PaymentStatusFailed {
		metrics.PaymentsProcessed.WithLabelValues("gateway", "failed").Inc()
		metrics.ActiveOrders.WithLabelValues("pending").Dec()
		logger.Logger.Info("payment callback completed",
			zap.String("status", string(verifyResp.Status)),
			zap.Duration("latency", time.Since(start)),
		)
		_ = s.messenger.SendToUser(ctx, int64(payment.UserID), msgerror.PaymentFailed) // ✅ FIX: payment.UserID
		return nil
	}

	if verifyResp.Status != paymententity.PaymentStatusSuccess {
		return nil
	}

	// 2. Update Order to PAID
	if ouErr := s.orderSvc.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
		OrderID: payment.OrderID,
		Status:  orderentity.OrderStatusPaid,
	}); ouErr != nil {
		return richerror.New(Op, ouErr).
			WithKind(richerror.KindQueryFailure).WithMessage(msgerror.OrderUpdateFailed)
	}

	// 3. Get order and fulfill
	order, obErr := s.orderSvc.GetById(ctx, payment.OrderID)
	if obErr != nil {
		metrics.PaymentsProcessed.WithLabelValues("gateway", "failed").Inc()
		metrics.CheckoutLatency.WithLabelValues("payment_callback").Observe(time.Since(start).Seconds())
		logger.Logger.Error("payment callback failed", zap.Error(vErr), zap.Uint64("order_id", payment.OrderID))
		return richerror.New(Op, obErr)
	}

	// 4. Fulfill (async with recover)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Log panic with Zap
			}
		}()
		s.fulfillOrderAsync(order)
	}()

	metrics.PaymentsProcessed.WithLabelValues("gateway", "success").Inc()
	metrics.ActiveOrders.WithLabelValues("paid").Dec()
	metrics.ActiveOrders.WithLabelValues("processing").Inc()
	logger.Logger.Info("payment callback completed",
		zap.String("status", string(verifyResp.Status)),
		zap.Duration("latency", time.Since(start)),
	)
	// 5. Notify
	_ = s.messenger.SendToUser(ctx, int64(payment.UserID),
		fmt.Sprintf("Payment successful! Order #%d is being processed.", order.ID))
	_ = s.messenger.SendToAdminChannel(ctx,
		fmt.Sprintf("New order #%d paid via %s", order.ID, payment.Method))

	return nil
}
