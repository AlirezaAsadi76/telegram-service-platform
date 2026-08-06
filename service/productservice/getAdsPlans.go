package productservice

import (
	"context"

	"telegram-service-platform/params"
	"telegram-service-platform/pkg/mapper"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) GetAdsPlans(ctx context.Context) (params.GetAdsPlansResponse, error) {
	const Op = "productservice.GetAdsPlans"

	plans, err := s.repository.GetAdsPlans(ctx)

	if err != nil {

		return params.GetAdsPlansResponse{},
			richerror.New(Op, err).
				WithKind(richerror.KindUnexpected).
				WithMessage(msgerror.Unexpected)
	}

	return mapper.MapAdsPlansResponse(plans), nil
}
