package productservice

import (
	"context"
	"telegram-service-platform/params"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) GetStarPlans(ctx context.Context) (params.GetStarPlansResponse, error) {

	const Op = "productservice.GetStarPlans"

	plans, err := s.repository.GetStarPlans(ctx)
	if err != nil {

		return params.GetStarPlansResponse{},
			richerror.New(Op, err)
	}

	priceStars, psErr := s.priceRepo.GetStarPrice(ctx)
	if psErr != nil {
		return params.GetStarPlansResponse{}, richerror.New(Op, psErr).WithKind(richerror.KindInfrastructure).WithMessage(msgerror.CacheReadFailed)
	}

	//TODO - add price to plans response

	response := params.GetStarPlansResponse{
		Plans: make([]params.StarPlanInfo, 0, len(plans)),
	}

	for _, plan := range plans {
		price, cErr := s.pricingSVc.CalculatePrice(ctx, float64(plan.Amount)*priceStars.PricePerStar)
		if cErr != nil {
			return response,
				richerror.New(Op, cErr).
					WithKind(richerror.KindUnexpected).
					WithMessage(msgerror.Unexpected)
		}

		response.Plans = append(response.Plans, params.StarPlanInfo{
			Price:  price,
			ID:     plan.ID,
			Amount: plan.Amount,
		})
	}

	return response, nil
}
