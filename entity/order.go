package entity

import (
	"telegram-service-platform/entity/productentity"
	"time"
)

type Order struct {
	ID uint64

	UserID uint64

	ProductType productentity.ProductType

	ProductID uint64

	Quantity int64

	Amount float64

	Currency Currency

	Status Status

	Metadata map[string]any

	CreatedAt time.Time

	UpdatedAt time.Time
}
