package smmproviderservice

import (
	"context"
	"fmt"
	"telegram-service-platform/entity/providerentity"
	"telegram-service-platform/params/smmparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) CreateOrder(ctx context.Context, req smmparams.CreateOrderAdapterRequest) (smmparams.CreateOrderResult, error) {
	const Op = "smmproviderservice.CreateOrder"

	providers, err := s.repo.GetActiveByType(ctx, providerentity.ProviderTypeSMM)

	if err != nil {
		return smmparams.CreateOrderResult{}, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	if len(providers) == 0 {
		return smmparams.CreateOrderResult{
			Outcome: smmparams.CreateOrderOutcomeUnknown,
		}, richerror.New(
			Op,
			fmt.Errorf("no active smm provider"),
		).
			WithKind(richerror.KindDependencyFailure).
			WithCode(richerror.CodeSMMProviderUnavailable).
			WithMessage(msgerror.NoAvailableAdapter)
	}

	rejected := 0
	usableProviders := 0

	for _, provider := range providers {
		adapter, ok := s.providers[provider.Name]
		if !ok {
			continue
		}

		breaker, okk := s.breakers[provider.Name]
		if !okk {
			continue
		}

		usableProviders++

		if !breaker.Allow() {
			continue
		}

		response, createErr := adapter.Create(ctx, req)

		if createErr != nil {
			if richerror.IsCode(createErr, richerror.CodeSMMProviderRejected) {
				rejected++
				continue
			}

			breaker.RecordFailure()

			return smmparams.CreateOrderResult{
				Outcome:      smmparams.CreateOrderOutcomeUnknown,
				ProviderID:   provider.ID,
				ProviderName: provider.Name,
			}, richerror.New(Op, createErr).
				WithKind(richerror.KindExternalAPI).
				WithMessage(msgerror.ExternalServiceFailed)
		}

		switch response.Outcome {
		case smmparams.CreateOrderOutcomeCreated:
			if response.ExternalOrderID == "" {
				breaker.RecordFailure()

				return smmparams.CreateOrderResult{
					Outcome:      smmparams.CreateOrderOutcomeUnknown,
					ProviderID:   provider.ID,
					ProviderName: provider.Name,
				}, richerror.New(Op,
					fmt.Errorf("provider %s returned empty external order id", provider.Name)).
					WithKind(richerror.KindExternalAPI).
					WithCode(richerror.CodeSMMProviderInvalidResponse).
					WithMessage(msgerror.SMMProviderInvalidResponse)
			}

			breaker.RecordSuccess()

			return smmparams.CreateOrderResult{
				Outcome:         smmparams.CreateOrderOutcomeCreated,
				ProviderID:      provider.ID,
				ProviderName:    provider.Name,
				ExternalOrderID: response.ExternalOrderID,
			}, nil

		case smmparams.CreateOrderOutcomeRejected:
			rejected++
			continue

		case smmparams.CreateOrderOutcomeUnknown:
			breaker.RecordFailure()

			return smmparams.CreateOrderResult{
				Outcome:      smmparams.CreateOrderOutcomeUnknown,
				ProviderID:   provider.ID,
				ProviderName: provider.Name,
			}, richerror.New(Op, fmt.Errorf("provider returned unknown create outcome")).
				WithKind(richerror.KindExternalAPI).
				WithCode(richerror.CodeSMMProviderInvalidResponse).
				WithMessage(msgerror.SMMProviderInvalidResponse)
		}
	}

	if usableProviders == 0 {
		return smmparams.CreateOrderResult{
			Outcome: smmparams.CreateOrderOutcomeUnknown,
		}, richerror.New(
			Op,
			fmt.Errorf("no usable smm provider"),
		).
			WithKind(richerror.KindDependencyFailure).
			WithCode(richerror.CodeSMMProviderUnavailable).
			WithMessage(msgerror.SMMProviderUnavailable)
	}

	if rejected == usableProviders {
		return smmparams.CreateOrderResult{
			Outcome: smmparams.CreateOrderOutcomeRejected,
		}, nil
	}

	return smmparams.CreateOrderResult{
		Outcome: smmparams.CreateOrderOutcomeUnknown,
	}, richerror.New(Op, fmt.Errorf("no usable smm provider")).
		WithKind(richerror.KindDependencyFailure).
		WithCode(richerror.CodeSMMProviderUnavailable).
		WithMessage(msgerror.SMMProviderUnavailable)
}
