package orderentity

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"
	"time"
)

type Order struct {
	ID              uint64
	UserID          uint64
	ProductType     productentity.ProductType
	ProductID       uint64
	Quantity        int64
	TargetLink      string // لینک کانال/پست برای SMM
	Amount          entity.Amount
	Currency        entity.Currency
	Status          OrderStatus // PENDING → PAID → PROCESSING → COMPLETED/FAILED/CANCELLED
	ExternalOrderID string      // شناسه سفارش در SMM Provider
	ProviderID      *uint64     // کدام پروایدر سفارش را گرفت
	Metadata        map[string]any
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
