package mapper

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/params/productparams"
)

func MapAdsPlansResponse(
	plans []entity.AdsPlan,
) productparams.GetAdsPlansResponse {

	result := make(
		[]productparams.AdsPlanInfo,
		0,
		len(plans),
	)

	for _, plan := range plans {

		result = append(
			result,
			productparams.AdsPlanInfo{
				ID:             plan.ID,
				Views:          plan.Views,
				CPM:            plan.CPM,
				DailyViewLimit: plan.DailyViewLimit,
			},
		)
	}

	return productparams.GetAdsPlansResponse{
		Plans: result,
	}
}
