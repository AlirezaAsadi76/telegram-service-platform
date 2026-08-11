// service/walletservice/credit.go
package walletservice

import (
	"context"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/params/walletparam"
)

// Credit performs an atomic credit operation with idempotency guarantee
func (s *Service) Credit(ctx context.Context, req walletparam.CreditRequest) (*walletparam.CreditResponse, error) {
	const Op = "walletservice.Credit"

	// 1. Check idempotency - have we already processed this?
	existingTx, gtErr := s.txRepo.GetTransactionByIdempotencyKey(ctx, req.IdempotencyKey)
	if gtErr == nil && existingTx != nil {
		// Already processed - return existing result
		wallet, _ := s.repo.GetByUserID(ctx, req.UserID)
		return &walletparam.CreditResponse{
			TransactionID: existingTx.ID,
			NewBalance:    wallet.Balance,
		}, nil
	}

	// 2. Set idempotency key in Redis (prevents concurrent processing)
	ok, ifErr := s.idempotencyRepo.SetIfNotExists(ctx, req.IdempotencyKey, "processing", 300)
	if ifErr != nil {
		return nil, richerror.New(Op, ifErr).WithKind(richerror.KindIdempotencyFailure).WithMessage(msgerror.IdempotencyAlreadyProcessing)
	}
	if !ok {
		// Another process is handling this
		return nil, richerror.New(Op, ifErr).WithKind(richerror.KindIdempotencyFailure).WithMessage(msgerror.IdempotencyAlreadyProcessing)
	}

	// 3. Start DB transaction and lock wallet
	wallet, gfErr := s.repo.GetForUpdate(ctx, req.UserID)
	if gfErr != nil {
		_ = s.idempotencyRepo.Delete(ctx, req.IdempotencyKey)
		return nil, richerror.New(Op, gfErr).WithKind(richerror.KindNotFound).WithMessage(msgerror.WalletNotFound)
	}

	// 4. Create pending transaction
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
		_ = s.idempotencyRepo.Delete(ctx, req.IdempotencyKey)
		return nil, richerror.New(Op, err).WithKind(richerror.KindInfrastructure)
	}

	// 5. Update wallet balance (optimistic locking)
	newBalance := wallet.Balance + req.Amount
	newVersion := wallet.Version + 1
	if err := s.repo.UpdateBalanceAtomic(ctx, wallet.ID, newBalance, newVersion); err != nil {
		_ = s.idempotencyRepo.Delete(ctx, req.IdempotencyKey)
		return nil, richerror.New(Op, err).WithKind(richerror.KindInfrastructure)
	}

	// 6. Mark transaction as completed
	if err := s.txRepo.UpdateTransactionStatus(ctx, tx.ID, walletentity.WalletTransactionStatusComplete); err != nil {
		// Log error but don't fail - balance is already updated
		// This is a rare case that needs manual reconciliation
	}

	// 7. Mark idempotency as completed
	_, _ = s.idempotencyRepo.SetIfNotExists(ctx, req.IdempotencyKey, "completed", 86400) // 24h retention

	return &walletparam.CreditResponse{
		TransactionID: tx.ID,
		NewBalance:    newBalance,
	}, nil
}
