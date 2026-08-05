package mapper

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/params"
)

func MapStarPlansResponse(plans []entity.StarPackage) params.GetStarPlansResponse {

	res := params.GetStarPlansResponse{}

	for _, plan := range plans {
		res.Plans = append(res.Plans, params.StarPlanInfo{
			ID:     plan.ID,
			Amount: plan.Amount,
		})
	}
	return res
}
