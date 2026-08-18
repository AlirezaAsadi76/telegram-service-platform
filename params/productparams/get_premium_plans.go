package productparams

import "telegram-service-platform/entity/productentity"

type PremiumPlanInfo struct {
	ID     uint64
	Months uint8
	Price  productentity.Price
}

type GetPremiumPlansResponse struct {
	Plans []PremiumPlanInfo
}
