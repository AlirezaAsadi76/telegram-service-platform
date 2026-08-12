package orderparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"
)

type CreateRequest struct {
	UserID      uint64
	ProductType productentity.ProductType // from productentity
	ProductID   uint64
	Quantity    int64
	TargetLink  string
	Amount      entity.Amount
	Currency    entity.Currency
}

type CreateResponse struct {
	OrderID uint64
}
