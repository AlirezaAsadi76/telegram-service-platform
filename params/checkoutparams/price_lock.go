package checkoutparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"
)

type LockSMMPriceRequest struct {
	UserID        uint64
	ProductType   productentity.ProductType
	ProductID     uint64
	Quantity      int64
	Currency      entity.Currency
	PaymentMethod entity.PriceLockPaymentMethod
}

type LockSMMPriceResponse struct {
	Amount        entity.Amount
	Currency      entity.Currency
	PaymentMethod entity.PriceLockPaymentMethod
	LockedAt      int64
	ExpiresAt     int64
}
