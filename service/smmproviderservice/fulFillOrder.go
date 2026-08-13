package smmproviderservice

import (
	"context"
	"fmt"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/providerentity"
	"telegram-service-platform/params/smmprams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) FulfillOrder(ctx context.Context, order *orderentity.Order) error {
	const Op = "smmproviderservice.FulfillOrder"
	dbProviders, err := s.repo.GetActiveByType(ctx, providerentity.ProviderTypeSMM)
	if err != nil {
		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).WithMessage(msgerror.ProviderQueryFailed)
	}

	for _, p := range dbProviders {
		adapter, ok := s.providers[p.Name]
		if !ok {
			continue
		}

		cb := s.breakers[p.Name]
		if !cb.Allow() {
			continue
		}

		req := smmprams.CreateOrderAdapterRequest{ // ✅ FIX: smmproviderservice.CreateOrderRequest
			ServiceID: fmt.Sprintf("%d", order.ProductID),
			Link:      order.TargetLink,
			Quantity:  order.Quantity,
		}
		resp, cErr := adapter.Create(ctx, req)
		if cErr != nil {
			cb.RecordFailure()
			continue
		}

		cb.RecordSuccess()
		order.ExternalOrderID = resp.ExternalOrderID
		order.Status = orderentity.OrderStatusProcessing
		return nil
	}

	return richerror.New(Op, fmt.Errorf("no available Adapter")).
		WithKind(richerror.KindExternalAPI).WithMessage(msgerror.NoAvailableAdapter)
}
