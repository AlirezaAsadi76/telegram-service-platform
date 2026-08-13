package paymententity

import (
	"telegram-service-platform/entity"
	"time"
)

type Payment struct {
	ID             uint64
	OrderID        uint64
	UserID         uint64
	Method         PaymentMethod // GATEWAY / CRYPTO
	Amount         entity.Amount
	Currency       entity.Currency
	Status         PaymentStatus // PENDING → PROCESSING → SUCCESS / FAILED / EXPIRED
	ExternalID     string        // شناسه درگاه پرداخت
	IdempotencyKey string
	CallbackData   map[string]any // داده‌های کال‌بک
	ExpiredAt      time.Time      // برای crypto
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
