package params

import "telegram-service-platform/entity/productentity"

type StarPlanInfo struct {
	ID     uint64
	Amount uint64
	Price  productentity.Price
}

type GetStarPlansResponse struct {
	Plans []StarPlanInfo
}
