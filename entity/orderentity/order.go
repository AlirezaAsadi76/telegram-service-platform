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
	TargetLink      string
	Amount          entity.Amount
	Currency        entity.Currency
	Status          OrderStatus
	ExternalOrderID string
	ProviderID      *uint64
	Metadata        map[string]any
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
