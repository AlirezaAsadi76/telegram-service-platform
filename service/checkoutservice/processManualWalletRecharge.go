package checkoutservice

import (
	"context"
	"fmt"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/walletparam"
	"telegram-service-platform/pkg/hashing"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/richerror"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) ProcessManualWalletRecharge(ctx context.Context, req checkoutparams.ManualRechargeRequest) error {
	const Op = "checkoutservice.ProcessManualWalletRecharge"

	start := time.Now()
	logger.Logger.Info("checkout manual recharge started",
		zap.Uint64("user_id", req.UserID),
		zap.Int64("amount", int64(req.Amount)),
	)

	nonce := uuid.New().String()
	idempotencyKey := hashing.EncodeStringToSHA256(
		fmt.Sprintf("%s:%d:%d:%d:%s", s.config.PrefixManualIdempotencyKey,
			req.AdminID, req.UserID, req.Amount, nonce))

	tx, err := s.walletSvc.Credit(ctx, walletparam.CreditRequest{
		UserID:         req.UserID,
		Amount:         req.Amount,
		ReferenceID:    fmt.Sprintf("manual_admin:%d", req.AdminID),
		IdempotencyKey: idempotencyKey,
	})

	if err != nil {
		metrics.WalletTransactions.WithLabelValues("credit_failed").Inc()
		metrics.CheckoutLatency.WithLabelValues("manual_recharge").Observe(time.Since(start).Seconds())
		logger.Logger.Error("checkout manual recharge failed", zap.Error(err), zap.Uint64("user_id", req.UserID),
			zap.Duration("latency", time.Since(start)))
		return richerror.New(Op, err)
	}

	metrics.WalletTransactions.WithLabelValues("credit").Inc()
	metrics.CheckoutLatency.WithLabelValues("manual_recharge").Observe(time.Since(start).Seconds())

	logger.Logger.Info("checkout manual recharge completed",
		zap.Uint64("wallet_tx_id", tx.TransactionID),
		zap.Duration("latency", time.Since(start)),
	)

	_ = s.messenger.SendToUser(ctx, int64(req.UserID),
		fmt.Sprintf("Your wallet charged: %d by admin", req.Amount))
	_ = s.messenger.SendToAdminChannel(ctx,
		fmt.Sprintf("Admin %d recharged user %d: %d", req.AdminID, req.UserID, req.Amount))

	return nil
}
