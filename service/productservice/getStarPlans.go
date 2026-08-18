package productservice

import (
	"context"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) GetStarPlans(ctx context.Context) (productparams.GetStarPlansResponse, error) {

	const Op = "productservice.GetStarPlans"

	plans, err := s.repository.GetStarPlans(ctx)
	if err != nil {

		return productparams.GetStarPlansResponse{},
			richerror.New(Op, err)
	}

	response := productparams.GetStarPlansResponse{
		Plans: make([]productparams.StarPlanInfo, 0, len(plans)),
	}

	for _, plan := range plans {
		price, cErr := s.pricingSVc.CalculateStarsPrice(ctx, float64(plan.Amount))
		if cErr != nil {
			return response,
				richerror.New(Op, cErr).
					WithKind(richerror.KindUnexpected).
					WithMessage(msgerror.Unexpected)
		}

		response.Plans = append(response.Plans, productparams.StarPlanInfo{
			Price:  price,
			ID:     plan.ID,
			Amount: plan.Amount,
		})
	}

	return response, nil
}
