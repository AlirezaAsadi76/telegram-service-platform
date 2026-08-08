package orderparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"
)

type CreateOrderRequest struct {
	UserID      uint64
	ProductType productentity.ProductType
	ProductID   uint64
	Quantity    int64
	Amount      float64
	Currency    entity.Currency
	Status      entity.Status
	Metadata    map[string]any
}

type CreateOrderResponse struct {
	Id     uint64
	Status entity.Status
}
