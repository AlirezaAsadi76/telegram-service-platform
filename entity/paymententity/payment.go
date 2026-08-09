package paymententity

import (
	"telegram-service-platform/entity"
	"time"
)

type Payment struct {
	ID        uint64
	OrderID   uint64
	UserID    uint64
	Method    PaymentMethod
	Provider  string
	Amount    float64
	Currency  entity.Currency
	Status    PaymentStatus
	Reference string
	Metadata  map[string]any
	CreatedAt time.Time
	UpdatedAt time.Time
}
