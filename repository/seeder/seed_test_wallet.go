package seeder

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/logger"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// SeedTestWallet کیف پول و تراکنش‌های تستی برای کاربر تستی ایجاد می‌کند
func (s *Seeder) SeedTestWallet(ctx context.Context) error {
	const op = "seeder.SeedTestWallet"

	user, err := s.userRepo.FindUserByTelegramID(ctx, testTelegramID)
	if err != nil {
		logger.Logger.Error("test user not found",
			zap.String("op", op),
			zap.Error(err),
		)
		return err
	}

	existingWallet, err := s.walletRepo.GetByUserID(ctx, user.ID)
	if err == nil && existingWallet != nil {
		logger.Logger.Info("test wallet already exists",
			zap.String("op", op),
			zap.Uint64("user_id", user.ID),
		)
		return nil
	}

	// ساخت کیف پول جدید
	wallet := &walletentity.Wallet{
		UserID:    user.ID,
		Balance:   entity.Amount(decimal.NewFromInt(100000)), // ۱۰۰,۰۰۰ تومان
		Currency:  "TOMAN",
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.walletRepo.Create(ctx, wallet)
	if err != nil {
		logger.Logger.Error("failed to create test wallet",
			zap.String("op", op),
			zap.Error(err),
		)
		return err
	}

	logger.Logger.Info("test wallet created successfully",
		zap.String("op", op),
		zap.Uint64("user_id", user.ID),
		zap.String("balance", wallet.Balance.String()),
	)

	// ساخت تراکنش‌های تستی
	return s.seedTestTransactions(ctx, wallet.ID, user.ID)
}

// seedTestTransactions تراکنش‌های تستی ایجاد می‌کند
func (s *Seeder) seedTestTransactions(ctx context.Context, walletID, userID uint64) error {
	const op = "seeder.seedTestTransactions"

	transactions := []walletentity.WalletTransaction{
		{
			WalletID:       walletID,
			UserID:         userID,
			Type:           walletentity.WalletTransactionTypeDeposit,
			Amount:         entity.Amount(decimal.NewFromInt(50000)),
			Status:         walletentity.WalletTransactionStatusComplete,
			ReferenceID:    "deposit_1",
			IdempotencyKey: "seed_tx_idempotency_1",
			CreatedAt:      time.Now().Add(-5 * 24 * time.Hour),
		},
		{
			WalletID:       walletID,
			UserID:         userID,
			Type:           walletentity.WalletTransactionTypeDeposit,
			Amount:         entity.Amount(decimal.NewFromInt(30000)),
			Status:         walletentity.WalletTransactionStatusComplete,
			ReferenceID:    "deposit_2",
			IdempotencyKey: "seed_tx_idempotency_2",
			CreatedAt:      time.Now().Add(-3 * 24 * time.Hour),
		},
		{
			WalletID:       walletID,
			UserID:         userID,
			Type:           walletentity.WalletTransactionTypeWithdraw,
			Amount:         entity.Amount(decimal.NewFromInt(10000)),
			Status:         walletentity.WalletTransactionStatusComplete,
			ReferenceID:    "withdraw_1",
			IdempotencyKey: "seed_tx_idempotency_3",
			CreatedAt:      time.Now().Add(-1 * 24 * time.Hour),
		},
	}

	for _, tx := range transactions {
		err := s.walletRepo.CreateTransaction(ctx, &tx)
		if err != nil {
			logger.Logger.Error("failed to create test transaction",
				zap.String("op", op),
				zap.Error(err),
			)
			return err
		}
	}

	logger.Logger.Info("test transactions created successfully",
		zap.String("op", op),
		zap.Int("count", len(transactions)),
	)

	return nil
}
