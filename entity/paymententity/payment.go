package paymententity

import (
	"telegram-service-platform/entity"
	"time"
)

type Payment struct {
	ID             uint64
	OrderID        uint64
	UserID         uint64
	Method         PaymentMethod
	Amount         entity.Amount
	Currency       entity.Currency
	Status         PaymentStatus
	ExternalID     string
	IdempotencyKey string
	CallbackData   map[string]any
	ExpiredAt      time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
