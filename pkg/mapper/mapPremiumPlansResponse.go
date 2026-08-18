package mapper

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/params/productparams"
)

func MapPremiumPlansResponse(plans []entity.PremiumPlan) productparams.GetPremiumPlansResponse {

	result := make([]productparams.PremiumPlanInfo, 0, len(plans))

	for _, plan := range plans {

		result = append(
			result,
			productparams.PremiumPlanInfo{
				ID:     plan.ID,
				Months: plan.Duration.Months(),
			},
		)
	}

	return productparams.GetPremiumPlansResponse{
		Plans: result,
	}
}
