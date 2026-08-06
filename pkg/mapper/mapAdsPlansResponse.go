package mapper

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/params"
)

func MapAdsPlansResponse(
	plans []entity.AdsPlan,
) params.GetAdsPlansResponse {

	result := make(
		[]params.AdsPlanInfo,
		0,
		len(plans),
	)

	for _, plan := range plans {

		result = append(
			result,
			params.AdsPlanInfo{
				ID:             plan.ID,
				Views:          plan.Views,
				CPM:            plan.CPM,
				DailyViewLimit: plan.DailyViewLimit,
			},
		)
	}

	return params.GetAdsPlansResponse{
		Plans: result,
	}
}
