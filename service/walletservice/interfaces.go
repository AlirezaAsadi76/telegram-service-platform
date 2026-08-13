package walletservice

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/walletentity"
	"time"
)

type Repository interface {
	GetByUserID(ctx context.Context, userID uint64) (*walletentity.Wallet, error)
	Create(ctx context.Context, wallet *walletentity.Wallet) error
	UpdateBalanceAtomic(ctx context.Context, walletID uint64, newBalance entity.Amount, newVersion int64) error
	GetForUpdate(ctx context.Context, userID uint64) (*walletentity.Wallet, error)
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *walletentity.WalletTransaction) error
	UpdateTransactionStatus(ctx context.Context, txID uint64, status walletentity.WalletTransactionStatus) error
	GetTransactionByIdempotencyKey(ctx context.Context, key string) (*walletentity.WalletTransaction, error)
}

type IdempotencyChecker interface {
	SetIfNotExists(ctx context.Context, idempotencyKey string, Value entity.IdempotencyStatus, ttl time.Duration) (bool, error)
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Set(ctx context.Context, idempotencyKey string, Value entity.IdempotencyStatus, ttl time.Duration) error
}
