package mapper

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/params"
)

func MapPremiumPlansResponse(plans []entity.PremiumPlan) params.GetPremiumPlansResponse {

	result := make([]params.PremiumPlanInfo, 0, len(plans))

	for _, plan := range plans {

		result = append(
			result,
			params.PremiumPlanInfo{
				ID:     plan.ID,
				Months: plan.Months,
			},
		)
	}

	return params.GetPremiumPlansResponse{
		Plans: result,
	}
}
