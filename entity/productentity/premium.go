package productentity

import "telegram-service-platform/entity"

type PremiumPrice struct {
	Months   uint8
	PriceUSD entity.Amount
}
