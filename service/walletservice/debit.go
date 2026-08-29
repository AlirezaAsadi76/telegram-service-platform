package walletservice

import (
	"context"
	"fmt"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/params/walletparam"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) Debit(ctx context.Context, req walletparam.DebitRequest) (*walletparam.DebitResponse, error) {
	const Op = "walletservice.Debit"

	// 1. Idempotency check
	existingTx, gtErr := s.txRepo.GetTransactionByIdempotencyKey(ctx, req.IdempotencyKey)
	if gtErr == nil && existingTx != nil {
		wallet, _ := s.repo.GetByUserID(ctx, req.UserID)
		return &walletparam.DebitResponse{
			TransactionID: existingTx.ID,
			NewBalance:    wallet.Balance,
		}, nil
	}
	// 2. Set idempotency key in Redis (prevents concurrent processing)
	ok, ifErr := s.idempotencyRepo.SetIfNotExists(ctx, req.IdempotencyKey, entity.IdempotencyStatusProcessing, s.config.IdempotencyProcessingTTL)
	if ifErr != nil {
		return nil, richerror.New(Op, ifErr).WithKind(richerror.KindIdempotencyFailure).WithMessage(msgerror.IdempotencyAlreadyProcessing)
	}
	if !ok {
		// Another process is handling this
		return nil, richerror.New(Op, ifErr).WithKind(richerror.KindIdempotencyFailure).WithMessage(msgerror.IdempotencyAlreadyProcessing)
	}
	// 2. Lock wallet
	wallet, gfErr := s.repo.GetForUpdate(ctx, req.UserID)
	if gfErr != nil {
		_ = s.idempotencyRepo.Delete(ctx, req.IdempotencyKey)
		return nil, richerror.New(Op, gfErr).WithKind(richerror.KindNotFound).WithMessage(msgerror.WalletNotFound)
	}

	// 3. Check sufficient balance
	if !wallet.HasSufficient(req.Amount) {
		_ = s.idempotencyRepo.Delete(ctx, req.IdempotencyKey)
		return nil, richerror.New(Op, fmt.Errorf("insufficient balance: have %d, need %d", wallet.Balance, req.Amount))
	}

	// 4. Create pending transaction
	tx := &walletentity.WalletTransaction{
		WalletID:       wallet.ID,
		UserID:         req.UserID,
		Type:           walletentity.WalletTransactionTypeWithdraw,
		Amount:         req.Amount,
		Status:         walletentity.WalletTransactionStatusPending,
		ReferenceID:    req.ReferenceID,
		IdempotencyKey: req.IdempotencyKey,
	}
	if err := s.txRepo.CreateTransaction(ctx, tx); err != nil {
		_ = s.idempotencyRepo.Delete(ctx, req.IdempotencyKey)
		return nil, richerror.New(Op, err).WithKind(richerror.KindInfrastructure)
	}

	// 5. Debit balance
	newBalance := wallet.Balance.Sub(req.Amount)
	newVersion := wallet.Version + 1
	if err := s.repo.UpdateBalanceAtomic(ctx, wallet.ID, newBalance, newVersion); err != nil {
		_ = s.idempotencyRepo.Delete(ctx, req.IdempotencyKey)
		return nil, fmt.Errorf("concurrent update: %w", err)
	}

	// 6. Mark transaction as completed
	if err := s.txRepo.UpdateTransactionStatus(ctx, tx.ID, walletentity.WalletTransactionStatusComplete); err != nil {
		// Log error but don't fail - balance is already updated
		// This is a rare case that needs manual reconciliation
	}

	// 7. Mark idempotency as completed
	_ = s.idempotencyRepo.Set(ctx, req.IdempotencyKey, entity.IdempotencyStatusComplete, s.config.IdempotencyCompletedTTL) // 24h retention

	return &walletparam.DebitResponse{
		TransactionID: tx.ID,
		NewBalance:    newBalance,
	}, nil
}
