package walletentity

import (
	"telegram-service-platform/entity"
	"time"
)

type Wallet struct {
	ID        uint64
	UserID    uint64
	Balance   entity.Amount
	Currency  entity.Currency
	Version   int64 // برای Optimistic Locking
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (w Wallet) HasSufficient(amount entity.Amount) bool {

	if w.Balance.Value() < amount.Value() {
		return false
	}
	return true
}
