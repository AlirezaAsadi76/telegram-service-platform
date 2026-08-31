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

	"github.com/go-telegram/bot"
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

	// ۲. اگر پرداخت ناموفق بود، اطلاع‌رسانی و خروج
	if verifyResp.Status == paymententity.PaymentStatusFailed || verifyResp.Status == paymententity.PaymentStatusCanceled {
		metrics.PaymentsProcessed.WithLabelValues("gateway", "failed").Inc()

		payment, _ := s.paymentSvc.GetById(ctx, paymentID)
		if payment != nil {
			_ = s.messenger.Send(ctx, &bot.SendMessageParams{
				ChatID: payment.UserID,
				Text:   "❌ پرداخت ناموفق بود. لطفاً دوباره تلاش کنید یا از روش دیگری استفاده کنید.",
			})
		}
		return nil
	}

	if verifyResp.Status != paymententity.PaymentStatusSuccess {
		return nil
	}

	payment, err := s.paymentSvc.GetById(ctx, paymentID)
	if err != nil {
		logger.Logger.Error("payment callback: failed to get payment details", zap.Error(err), zap.Uint64("payment_id", paymentID))
		return richerror.New(Op, err)
	}

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
		logger.Logger.Error("payment callback failed", zap.Error(obErr), zap.Uint64("order_id", payment.OrderID))
		return richerror.New(Op, obErr)
	}

	// 4. Fulfill (async with recover)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Logger.Error("panic in fulfillOrderAsync",
					zap.Uint64("order_id", order.ID),
					zap.Any("panic", r))
			}
		}()
		s.fulfillOrderAsync(order)
	}()

	metrics.PaymentsProcessed.WithLabelValues("gateway", "success").Inc()
	metrics.ActiveOrders.WithLabelValues("paid").Inc()
	metrics.ActiveOrders.WithLabelValues("processing").Inc()
	logger.Logger.Info("payment callback completed",
		zap.String("status", string(verifyResp.Status)),
		zap.Duration("latency", time.Since(start)),
	)
	// 5. Notify

	_ = s.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID: payment.UserID,
		Text:   fmt.Sprintf("Payment successful! Order #%d is being processed.", order.ID),
	})
	//TODO - send to admin
	//_ = s.messenger.Send(ctx, &bot.SendMessageParams{
	//	ChatID: payment.,
	//	Text:  fmt.Sprintf("New order #%d paid via %s", order.ID, payment.Method),
	//})

	metrics.CheckoutLatency.WithLabelValues("payment_callback").Observe(time.Since(start).Seconds())
	return nil
}
