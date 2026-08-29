package productparams

import (
	"telegram-service-platform/entity/productentity"
)

type StarPlanInfo struct {
	ID     uint64
	Amount int64
	Price  productentity.Price
}

type GetStarPlansResponse struct {
	Plans []StarPlanInfo
}
