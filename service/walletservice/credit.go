package walletservice

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/params/walletparam"

	"go.uber.org/zap"
)

// Credit performs an atomic credit operation with idempotency guarantee
func (s *Service) Credit(ctx context.Context, req walletparam.CreditRequest) (*walletparam.CreditResponse, error) {
	const Op = "walletservice.Credit"
	logger.Logger.Debug("check error manual wallet ",
		zap.Any("req", req))
	existingTx, gtErr := s.txRepo.GetTransactionByIdempotencyKey(ctx, req.IdempotencyKey)
	if gtErr == nil && existingTx != nil {
		wallet, _ := s.repo.GetByUserID(ctx, req.UserID)
		return &walletparam.CreditResponse{
			TransactionID: existingTx.ID,
			NewBalance:    wallet.Balance,
		}, nil
	}
	logger.Logger.Debug("check error manual wallet ",
		zap.Bool("existingTx", existingTx != nil),
		zap.Error(gtErr))

	ok, ifErr := s.idempotencyRepo.SetIfNotExists(ctx, req.IdempotencyKey, entity.IdempotencyStatusProcessing, s.config.IdempotencyProcessingTTL)
	if ifErr != nil {
		return nil, richerror.New(Op, ifErr).WithKind(richerror.KindIdempotencyFailure).WithMessage(msgerror.IdempotencyAlreadyProcessing)
	}
	if !ok {
		return nil, richerror.New(Op, ifErr).WithKind(richerror.KindIdempotencyFailure).WithMessage(msgerror.IdempotencyAlreadyProcessing)
	}
	logger.Logger.Debug("check error manual wallet ",
		zap.Bool("idempotency", ok),
		zap.Error(ifErr))

	success := false
	defer func() {
		if !success {
			_ = s.idempotencyRepo.Delete(ctx, req.IdempotencyKey)
		}
	}()

	wallet, gfErr := s.repo.GetForUpdate(ctx, req.UserID)
	if gfErr != nil {
		if delErr := s.idempotencyRepo.Delete(ctx, req.IdempotencyKey); delErr != nil {
			logger.Logger.Warn("failed to cleanup idempotency key after GetForUpdate failure",
				zap.String("op", Op),
				zap.String("idempotency_key", req.IdempotencyKey),
				zap.Error(delErr),
			)
		}
		return nil, richerror.New(Op, gfErr).WithKind(richerror.KindNotFound).WithMessage(msgerror.WalletNotFound)
	}

	tx := &walletentity.WalletTransaction{
		WalletID:       wallet.ID,
		UserID:         req.UserID,
		Type:           walletentity.WalletTransactionTypeDeposit,
		Amount:         req.Amount,
		Status:         walletentity.WalletTransactionStatusPending,
		ReferenceID:    req.ReferenceID,
		IdempotencyKey: req.IdempotencyKey,
	}
	if err := s.txRepo.CreateTransaction(ctx, tx); err != nil {
		if delErr := s.idempotencyRepo.Delete(ctx, req.IdempotencyKey); delErr != nil {
			logger.Logger.Warn("failed to cleanup idempotency key after CreateTransaction failure",
				zap.String("op", Op),
				zap.String("idempotency_key", req.IdempotencyKey),
				zap.Error(delErr),
			)
		}
		return nil, richerror.New(Op, err).WithKind(richerror.KindInfrastructure)
	}

	newBalance := wallet.Balance + req.Amount
	newVersion := wallet.Version + 1
	if err := s.repo.UpdateBalanceAtomic(ctx, wallet.ID, newBalance, newVersion); err != nil {
		if delErr := s.idempotencyRepo.Delete(ctx, req.IdempotencyKey); delErr != nil {
			logger.Logger.Warn("failed to cleanup idempotency key after UpdateBalanceAtomic failure",
				zap.String("op", Op),
				zap.String("idempotency_key", req.IdempotencyKey),
				zap.Error(delErr),
			)
		}
		return nil, richerror.New(Op, err).WithKind(richerror.KindInfrastructure)
	}

	// 6. Mark transaction as completed
	if err := s.txRepo.UpdateTransactionStatus(ctx, tx.ID, walletentity.WalletTransactionStatusComplete); err != nil {
		logger.Logger.Error("failed to update transaction status to COMPLETE",
			zap.String("op", Op),
			zap.Uint64("transaction_id", tx.ID),
			zap.Uint64("user_id", req.UserID),
			zap.Uint64("wallet_id", wallet.ID),
			zap.Any("amount", req.Amount),
			zap.Error(err),
		)
	}

	// 7. Mark idempotency as completed
	_ = s.idempotencyRepo.Set(ctx, req.IdempotencyKey, entity.IdempotencyStatusComplete, s.config.IdempotencyCompletedTTL) // 24h retention

	success = true
	return &walletparam.CreditResponse{
		TransactionID: tx.ID,
		NewBalance:    newBalance,
	}, nil
}
