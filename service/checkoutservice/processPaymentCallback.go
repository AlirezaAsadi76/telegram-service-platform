package checkoutservice

import (
	"context"
	"fmt"
	"time"

	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/richerror"

	"github.com/go-telegram/bot"
	"go.uber.org/zap"
)

func (s *Service) ProcessPaymentCallback(ctx context.Context, paymentID uint64, externalID string, callbackData map[string]any) error {
	const Op = "checkoutservice.ProcessPaymentCallback"

	start := time.Now()

	logger.Logger.Info(
		"payment callback started",
		zap.Uint64("payment_id", paymentID),
	)

	verifyResp, vErr := s.paymentSvc.ConfirmPayment(
		ctx,
		paymentparams.ConfirmPaymentRequest{
			PaymentID:    paymentID,
			ExternalID:   externalID,
			CallbackData: callbackData,
		},
	)
	if vErr != nil {
		metrics.PaymentsProcessed.
			WithLabelValues("gateway", "failed").
			Inc()

		metrics.CheckoutLatency.
			WithLabelValues("payment_callback").
			Observe(time.Since(start).Seconds())

		logger.Logger.Error(
			"payment callback failed",
			zap.Error(vErr),
			zap.Uint64("payment_id", paymentID),
		)

		return richerror.New(Op, vErr)
	}

	if verifyResp.Status == paymententity.PaymentStatusFailed ||
		verifyResp.Status == paymententity.PaymentStatusCanceled {
		metrics.PaymentsProcessed.
			WithLabelValues("gateway", "failed").
			Inc()

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
		logger.Logger.Error(
			"payment callback: failed to get payment details",
			zap.Error(err),
			zap.Uint64("payment_id", paymentID),
		)

		return richerror.New(Op, err)
	}

	metrics.PaymentsProcessed.
		WithLabelValues("gateway", "success").
		Inc()

	metrics.ActiveOrders.
		WithLabelValues("paid").
		Inc()

	logger.Logger.Info(
		"payment callback completed",
		zap.Uint64("payment_id", paymentID),
		zap.Uint64("order_id", payment.OrderID),
		zap.String("status", string(verifyResp.Status)),
		zap.Duration("latency", time.Since(start)),
	)

	_ = s.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID: payment.UserID,
		Text: fmt.Sprintf(
			"Payment successful! Order #%d is being processed.",
			payment.OrderID,
		),
	})

	metrics.CheckoutLatency.
		WithLabelValues("payment_callback").
		Observe(time.Since(start).Seconds())

	return nil
}
