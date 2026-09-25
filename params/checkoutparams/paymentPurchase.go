package checkoutparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/entity/productentity"
)

type DirectPaymentPurchase struct {
	UserID                 uint64
	ProductType            productentity.ProductType
	ProductID              uint64
	Quantity               int64
	TargetLink             string
	Amount                 entity.Amount
	Currency               entity.Currency
	Method                 paymententity.PaymentMethod
	PriceLockedAt          int64
	PriceLockExpiresAt     int64
	PriceLockPaymentMethod entity.PriceLockPaymentMethod
}

type PaymentURLResponse struct {
	OrderID    uint64
	PaymentID  uint64
	PaymentURL string
	Amount     entity.Amount
	Currency   entity.Currency
}
