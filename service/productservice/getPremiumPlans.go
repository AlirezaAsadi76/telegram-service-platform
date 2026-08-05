package productservice

import (
	"context"

	"telegram-service-platform/params"
	"telegram-service-platform/pkg/mapper"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) GetPremiumPlans(ctx context.Context) (params.GetPremiumPlansResponse, error) {
	const Op = "productservice.GetPremiumPlans"

	plans, err := s.repository.GetPremiumPlans(ctx)

	if err != nil {
		return params.GetPremiumPlansResponse{},
			richerror.New(Op, err).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}

	return mapper.MapPremiumPlansResponse(plans), nil
}
