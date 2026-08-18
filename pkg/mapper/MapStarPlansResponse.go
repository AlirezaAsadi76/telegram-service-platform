package mapper

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/params/productparams"
)

func MapStarPlansResponse(plans []entity.StarPackage) productparams.GetStarPlansResponse {

	res := productparams.GetStarPlansResponse{}

	for _, plan := range plans {
		res.Plans = append(res.Plans, productparams.StarPlanInfo{
			ID:     plan.ID,
			Amount: plan.Amount,
		})
	}
	return res
}
