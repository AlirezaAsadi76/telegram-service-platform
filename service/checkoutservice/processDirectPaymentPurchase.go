package checkoutservice

import (
	"context"
	"fmt"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/hashing"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/pkg/ts"
	"time"

	"go.uber.org/zap"
)

func (s *Service) ProcessDirectPaymentPurchase(ctx context.Context, req checkoutparams.DirectPaymentPurchase) (*checkoutparams.PaymentURLResponse, error) {
	const Op = "checkoutservice.ProcessDirectPaymentPurchase"

	start := time.Now()
	logger.Logger.Info("checkout direct payment started",
		zap.Uint64("user_id", req.UserID),
		zap.Int64("amount", int64(req.Amount)),
	)

	// 1. Create Order (PENDING)
	orderResp, err := s.orderSvc.Create(ctx, orderparams.CreateRequest{
		UserID:      req.UserID,
		ProductType: req.ProductType,
		ProductID:   req.ProductID,
		Quantity:    req.Quantity,
		TargetLink:  req.TargetLink,
		Amount:      req.Amount,
		Currency:    req.Currency,
	})
	if err != nil {
		metrics.OrdersTotal.WithLabelValues("direct_payment", "failed").Inc()
		metrics.CheckoutLatency.WithLabelValues("direct_payment").Observe(time.Since(start).Seconds())
		logger.Logger.Error("checkout direct payment failed", zap.Error(err), zap.Uint64("user_id", req.UserID))
		return nil, richerror.New(Op, err)
	}

	// 2. Generate idempotency key
	idempotencyKey := hashing.EncodeStringToSHA256(
		fmt.Sprintf("%s:%d:%d:%d", s.config.PrefixDirectIdempotencyKey, req.UserID, orderResp.OrderID, ts.Now()),
	)

	// 3. Create Payment
	paymentResp, cpErr := s.paymentSvc.Create(ctx, paymentparams.CreateRequest{
		OrderID:        orderResp.OrderID,
		UserID:         req.UserID,
		Method:         req.Method,
		Amount:         req.Amount,
		Currency:       req.Currency,
		IdempotencyKey: idempotencyKey,
	})
	if cpErr != nil {
		metrics.OrdersTotal.WithLabelValues("direct_payment", "failed").Inc()
		metrics.CheckoutLatency.WithLabelValues("direct_payment").Observe(time.Since(start).Seconds())
		logger.Logger.Error("checkout direct payment failed", zap.Error(cpErr), zap.Uint64("order_id", orderResp.OrderID))
		return nil, richerror.New(Op, cpErr)
	}

	metrics.OrdersTotal.WithLabelValues("direct_payment", "pending").Inc()
	metrics.PaymentsProcessed.WithLabelValues(string(req.Method), "pending").Inc()
	metrics.ActiveOrders.WithLabelValues("pending").Inc()
	metrics.CheckoutLatency.WithLabelValues("direct_payment").Observe(time.Since(start).Seconds())

	logger.Logger.Info("checkout direct payment completed",
		zap.Uint64("order_id", orderResp.OrderID),
		zap.Uint64("payment_id", paymentResp.PaymentID),
		zap.Duration("latency", time.Since(start)),
	)

	return &checkoutparams.PaymentURLResponse{
		OrderID:    orderResp.OrderID,
		PaymentID:  paymentResp.PaymentID,
		PaymentURL: paymentResp.PaymentURL,
	}, nil
}
