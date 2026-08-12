package walletentity

import (
	"telegram-service-platform/entity"
	"time"
)

type WalletTransaction struct {
	ID             uint64
	WalletID       uint64
	UserID         uint64
	Type           WalletTransactionType // DEPOSIT / WITHDRAW / REFUND
	Amount         entity.Amount
	Status         WalletTransactionStatus // PENDING / COMPLETED / FAILED / REVERSED
	ReferenceID    string                  // payment_id یا order_id
	IdempotencyKey string                  // کلید یکتای idempotency
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
