package orderparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"
	"time"

	"github.com/shopspring/decimal"
)

type GetOrderByIdRequest struct {
	Id uint64
}
type GetOrderByIdResponse struct {
	OrderInfo OrderInfo
}

type OrderInfo struct {
	Id          uint64
	UserId      uint64
	ProductId   uint64
	ProductType productentity.ProductType
	Quantity    uint64
	Amount      decimal.Decimal
	Currency    string
	Status      entity.Status
	Metadata    map[string]any
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
