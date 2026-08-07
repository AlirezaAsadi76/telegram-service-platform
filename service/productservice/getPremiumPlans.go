package productservice

import (
	"context"
	"errors"
	"telegram-service-platform/entity/productentity"

	"telegram-service-platform/params"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) GetPremiumPlans(ctx context.Context) (params.GetPremiumPlansResponse, error) {
	const Op = "productservice.GetPremiumPlans"

	plans, err := s.repository.GetPremiumPlans(ctx)

	if err != nil {
		return params.GetPremiumPlansResponse{},
			richerror.New(Op, err)
	}

	pricePlans, pErr := s.priceRepo.GetPremiumPrices(ctx)
	if pErr != nil {
		return params.GetPremiumPlansResponse{}, richerror.New(Op, pErr).WithKind(richerror.KindInfrastructure).WithMessage(msgerror.CacheReadFailed)
	}

	priceMap := makePremiumPriceMap(pricePlans)
	response := params.GetPremiumPlansResponse{
		Plans: make([]params.PremiumPlanInfo, 0, len(plans)),
	}
	for _, plan := range plans {
		price, ok := priceMap[plan.Duration.Months()]
		if !ok {
			return response,
				richerror.New(Op, errors.New(msgerror.PremiumPriceNotFound)).
					WithKind(richerror.KindNotFound).
					WithMessage(msgerror.PremiumPriceNotFound)
		}

		calculatedPrice, cErr := s.pricingSVc.CalculatePrice(ctx, price.PriceUSD)
		if cErr != nil {
			return response,
				richerror.New(Op, cErr).
					WithKind(richerror.KindUnexpected).
					WithMessage(msgerror.Unexpected)
		}
		response.Plans = append(response.Plans, params.PremiumPlanInfo{
			ID:     plan.ID,
			Months: plan.Duration.Months(),
			Price:  calculatedPrice,
		})

	}

	return response, nil
}

func makePremiumPriceMap(prices []productentity.PremiumPrice) map[uint8]productentity.PremiumPrice {
	result := make(map[uint8]productentity.PremiumPrice, len(prices))

	for _, price := range prices {
		result[price.Months] = price
	}

	return result
}
