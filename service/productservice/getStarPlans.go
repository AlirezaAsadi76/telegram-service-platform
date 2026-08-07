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

	response := params.GetStarPlansResponse{
		Plans: make([]params.StarPlanInfo, 0, len(plans)),
	}

	for _, plan := range plans {
		price, cErr := s.pricingSVc.CalculateStarsPrice(ctx, float64(plan.Amount))
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
