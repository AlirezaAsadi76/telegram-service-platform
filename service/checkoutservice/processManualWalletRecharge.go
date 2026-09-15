package checkoutservice

import (
	"context"
	"fmt"
	"time"

	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/params/walletparam"
	"telegram-service-platform/pkg/hashing"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/richerror"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) ProcessManualWalletRecharge(ctx context.Context, req checkoutparams.ManualRechargeRequest) error {
	const Op = "checkoutservice.ProcessManualWalletRecharge"

	start := time.Now()

	nonce := uuid.New().String()

	idempotencyKey := hashing.EncodeStringToSHA256(
		fmt.Sprintf(
			"%s:%d:%d:%d:%s",
			s.config.PrefixManualIdempotencyKey,
			req.AdminID,
			req.UserID,
			req.Amount,
			nonce,
		),
	)

	logger.Logger.Info(
		"checkout manual recharge started",
		zap.Uint64("user_id", req.UserID),
		zap.String("amount", req.Amount.String()),
		zap.String("nonce", nonce),
		zap.String("idempotency_key", idempotencyKey[:16]+"..."),
	)

	tx, err := s.walletSvc.Credit(
		ctx,
		walletparam.CreditRequest{
			UserID:         req.UserID,
			Amount:         req.Amount,
			ReferenceID:    fmt.Sprintf("manual_admin:%d", req.AdminID),
			IdempotencyKey: idempotencyKey,
		},
	)
	if err != nil {
		metrics.WalletTransactions.
			WithLabelValues("credit_failed").
			Inc()

		metrics.CheckoutLatency.
			WithLabelValues("manual_recharge").
			Observe(time.Since(start).Seconds())

		logger.Logger.Error(
			"checkout manual recharge failed",
			zap.Error(err),
			zap.Uint64("user_id", req.UserID),
			zap.String("idempotency_key", idempotencyKey[:16]+"..."),
			zap.Duration("latency", time.Since(start)),
		)

		return richerror.New(Op, err)
	}

	metrics.WalletTransactions.
		WithLabelValues("credit").
		Inc()

	metrics.CheckoutLatency.
		WithLabelValues("manual_recharge").
		Observe(time.Since(start).Seconds())

	logger.Logger.Info(
		"checkout manual recharge completed",
		zap.Uint64("wallet_transaction_id", tx.TransactionID),
		zap.String("idempotency_key", idempotencyKey[:16]+"..."),
		zap.Duration("latency", time.Since(start)),
	)

	if err := s.notificationSvc.Create(
		ctx,
		notificationparams.CreateRequest{
			UserID: req.UserID,
			Type:   notificationentity.NotificationTypeWalletRecharged,
			Payload: map[string]any{
				"amount":    req.Amount,
				"admin_id":  req.AdminID,
				"wallet_tx": tx.TransactionID,
			},
		},
	); err != nil {
		logger.Logger.Error(
			"user wallet recharge notification creation failed",
			zap.Uint64("user_id", req.UserID),
			zap.Uint64("wallet_transaction_id", tx.TransactionID),
			zap.Error(err),
		)
	}

	if err := s.notificationSvc.Create(
		ctx,
		notificationparams.CreateRequest{
			UserID: uint64(req.AdminID),
			Type:   notificationentity.NotificationTypeAdminWalletRecharge,
			Payload: map[string]any{
				"user_id":   req.UserID,
				"amount":    req.Amount,
				"wallet_tx": tx.TransactionID,
			},
		},
	); err != nil {
		logger.Logger.Error(
			"admin wallet recharge notification creation failed",
			zap.Uint64("admin_id", uint64(req.AdminID)),
			zap.Uint64("user_id", req.UserID),
			zap.Uint64("wallet_transaction_id", tx.TransactionID),
			zap.Error(err),
		)
	}

	return nil
}
