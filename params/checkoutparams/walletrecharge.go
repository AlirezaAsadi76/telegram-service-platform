package checkoutparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"
)

type WalletPurchaseRequest struct {
	UserID                 uint64
	ProductType            productentity.ProductType
	ProductID              uint64
	Quantity               int64
	TargetLink             string
	Amount                 entity.Amount
	Currency               entity.Currency
	IdempotencyKey         string
	PriceLockedAt          int64
	PriceLockExpiresAt     int64
	PriceLockPaymentMethod entity.PriceLockPaymentMethod
}

type WalletPurchaseResponse struct {
	OrderID    uint64
	WalletTxID uint64
	Amount     entity.Amount
	Currency   entity.Currency
	NewBalance entity.Amount
}

type ManualRechargeRequest struct {
	AdminID  entity.TelegramId
	UserID   uint64
	Amount   entity.Amount
	Currency entity.Currency
}
