package productservice

import (
	"context"
	"telegram-service-platform/params"
	"telegram-service-platform/pkg/mapper"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) GetStarPlans(ctx context.Context) (params.GetStarPlansResponse, error) {

	const Op = "productservice.GetStarPlans"

	plans, err := s.repository.GetStarPlans(ctx)
	if err != nil {

		return params.GetStarPlansResponse{},
			richerror.New(
				Op,
				err,
			).
				WithKind(
					richerror.KindUnexpected,
				).
				WithMessage(
					msgerror.Unexpected,
				)
	}
	//TODO - add price to plans response
	return mapper.MapStarPlansResponse(plans), nil
}
