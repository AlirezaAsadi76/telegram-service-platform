package walletparam

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/productentity"
)

type WalletPurchaseRequest struct {
	UserID         uint64
	ProductType    productentity.ProductType
	ProductID      uint64
	Quantity       int64
	TargetLink     string
	Amount         entity.Amount
	Currency       entity.Currency
	ReferenceID    string
	IdempotencyKey string
}

type WalletPurchaseResult struct {
	OrderID     uint64
	WalletTxID  uint64
	NewBalance  entity.Amount
	OrderStatus orderentity.OrderStatus
}
