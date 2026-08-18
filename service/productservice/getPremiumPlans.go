package productservice

import (
	"context"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) GetPremiumPlans(ctx context.Context) (productparams.GetPremiumPlansResponse, error) {
	const Op = "productservice.GetPremiumPlans"

	plans, err := s.repository.GetPremiumPlans(ctx)

	if err != nil {
		return productparams.GetPremiumPlansResponse{},
			richerror.New(Op, err)
	}

	response := productparams.GetPremiumPlansResponse{
		Plans: make([]productparams.PremiumPlanInfo, 0, len(plans)),
	}
	calculatedPrice, cErr := s.pricingSVc.CalculatePremiumPrices(ctx)
	if cErr != nil {
		return response,
			richerror.New(Op, cErr)
	}

	for _, plan := range plans {

		response.Plans = append(response.Plans, productparams.PremiumPlanInfo{
			ID:     plan.ID,
			Months: plan.Duration.Months(),
			Price:  calculatedPrice[plan.Duration.Months()],
		})

	}

	return response, nil
}
