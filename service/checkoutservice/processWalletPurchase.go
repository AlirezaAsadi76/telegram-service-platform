package checkoutservice

import (
	"context"
	"fmt"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/params/walletparam"
	"telegram-service-platform/pkg/hashing"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/pkg/ts"
	"time"

	"go.uber.org/zap"
)

func (s *Service) ProcessWalletPurchase(ctx context.Context, req checkoutparams.WalletPurchaseRequest) error {

	const Op = "checkoutservice.ProcessWalletPurchase"

	start := time.Now()
	logger.Logger.Info("checkout wallet purchase started",
		zap.Uint64("user_id", req.UserID),
		zap.Int64("amount", int64(req.Amount)),
	)
	// 1. Check balance
	balanceResp, err := s.walletSvc.GetBalance(ctx, walletparam.GetBalanceRequest{UserID: req.UserID})

	if err != nil {
		metrics.WalletTransactions.WithLabelValues("debit_failed").Inc()
		metrics.CheckoutLatency.WithLabelValues("wallet").Observe(time.Since(start).Seconds())
		logger.Logger.Error("checkout wallet purchase failed", zap.Error(err), zap.Uint64("user_id", req.UserID),
			zap.Duration("latency", time.Since(start)))

		return richerror.New(Op, err)

	}
	if balanceResp.Balance < req.Amount {
		metrics.WalletTransactions.WithLabelValues("debit_failed").Inc()
		metrics.CheckoutLatency.WithLabelValues("wallet").Observe(time.Since(start).Seconds())
		logger.Logger.Error("checkout wallet purchase failed", zap.Error(fmt.Errorf("insufficient balance")), zap.Uint64("user_id", req.UserID),
			zap.Int64("balance", int64(balanceResp.Balance)),
			zap.Int64("request_amount", int64(req.Amount)),
			zap.Duration("latency", time.Since(start)))

		return richerror.New(Op, fmt.Errorf("insufficient balance")).
			WithKind(richerror.KindValidation).WithMessage(msgerror.InsufficientBalance)
	}

	// 2. Create Order (PENDING)
	orderResp, oErr := s.orderSvc.Create(ctx, orderparams.CreateRequest{
		UserID:      req.UserID,
		ProductType: req.ProductType,
		ProductID:   req.ProductID,
		Quantity:    req.Quantity,
		TargetLink:  req.TargetLink,
		Amount:      req.Amount,
		Currency:    req.Currency,
	})

	if oErr != nil {
		metrics.WalletTransactions.WithLabelValues("debit_failed").Inc()
		metrics.CheckoutLatency.WithLabelValues("wallet").Observe(time.Since(start).Seconds())
		logger.Logger.Error("checkout wallet purchase failed", zap.Error(oErr), zap.Uint64("user_id", req.UserID),
			zap.Duration("latency", time.Since(start)))

		return richerror.New(Op, oErr)

	}

	// 3. Generate idempotency key

	idempotencyKey := hashing.EncodeStringToSHA256(
		fmt.Sprintf("%s:%d:%d:%d", s.config.PrefixWalletIdempotencyKey, req.UserID, orderResp.OrderID, ts.Now()))

	// 4. Debit Wallet
	_, debErr := s.walletSvc.Debit(ctx, walletparam.DebitRequest{
		UserID:         req.UserID,
		Amount:         req.Amount,
		ReferenceID:    fmt.Sprintf("order:%d", orderResp.OrderID),
		IdempotencyKey: idempotencyKey,
	})
	if debErr != nil {
		// Cancel order
		_ = s.orderSvc.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
			OrderID: orderResp.OrderID,
			Status:  orderentity.OrderStatusCanceled,
		})
		metrics.WalletTransactions.WithLabelValues("debit_failed").Inc()
		metrics.CheckoutLatency.WithLabelValues("wallet").Observe(time.Since(start).Seconds())
		logger.Logger.Error("checkout wallet purchase failed", zap.Error(debErr), zap.Uint64("order_id", orderResp.OrderID),
			zap.Duration("latency", time.Since(start)))
		return richerror.New(Op, debErr)
	}

	// 5. Update Order to PAID ✅ FIX: error handling
	if ouErr := s.orderSvc.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
		OrderID: orderResp.OrderID,
		Status:  orderentity.OrderStatusPaid,
	}); ouErr != nil {
		metrics.WalletTransactions.WithLabelValues("debit_failed").Inc()
		metrics.CheckoutLatency.WithLabelValues("wallet").Observe(time.Since(start).Seconds())
		logger.Logger.Error("checkout wallet purchase failed", zap.Error(ouErr), zap.Uint64("order_id", orderResp.OrderID),
			zap.Duration("latency", time.Since(start)))
		return richerror.New(Op, ouErr).
			WithKind(richerror.KindQueryFailure).WithMessage(msgerror.OrderUpdateFailed)
	}

	// 6. Get order and fulfill
	order, ogErr := s.orderSvc.GetById(ctx, orderResp.OrderID)
	if ogErr != nil {
		metrics.WalletTransactions.WithLabelValues("debit_failed").Inc()
		metrics.CheckoutLatency.WithLabelValues("wallet").Observe(time.Since(start).Seconds())
		logger.Logger.Error("checkout wallet purchase failed", zap.Error(ogErr), zap.Uint64("order_id", orderResp.OrderID),
			zap.Duration("latency", time.Since(start)))
		return richerror.New(Op, ogErr)
	}

	// 7. Fulfill (async with recover)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Log panic
			}
		}()
		s.fulfillOrderAsync(order)
	}()

	metrics.WalletTransactions.WithLabelValues("WalletPurchase").Inc()
	metrics.OrdersCreated.WithLabelValues("wallet", "processing").Inc()
	metrics.ActiveOrders.WithLabelValues("processing").Inc()
	metrics.CheckoutLatency.WithLabelValues("wallet").Observe(time.Since(start).Seconds())

	logger.Logger.Info("checkout wallet purchase completed",
		zap.Uint64("order_id", orderResp.OrderID),
		zap.Duration("latency", time.Since(start)),
	)
	// 8. Notify
	_ = s.messenger.SendToUser(ctx, int64(req.UserID),
		fmt.Sprintf("Order #%d placed! Processing...", order.ID))
	_ = s.messenger.SendToAdminChannel(ctx,
		fmt.Sprintf("New wallet order #%d from user %d", order.ID, req.UserID))

	return nil
}
