package checkoutparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"
)

type WalletPurchaseRequest struct {
	UserID      uint64
	ProductType productentity.ProductType
	ProductID   uint64
	Quantity    int64
	TargetLink  string
	Amount      entity.Amount
	Currency    entity.Currency
}

type ManualRechargeRequest struct {
	AdminID        entity.TelegramId
	UserTelegramID entity.TelegramId
	UserID         uint64
	Amount         entity.Amount
	Currency       entity.Currency
}
