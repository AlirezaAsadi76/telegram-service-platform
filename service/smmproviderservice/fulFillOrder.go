package smmproviderservice

import (
	"context"
	"fmt"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/providerentity"
	"telegram-service-platform/params/smmprams"
	"telegram-service-platform/pkg/richerror"
)

// FulfillOrder picks the best available provider and creates external order
func (s *Service) FulfillOrder(ctx context.Context, order *orderentity.Order) error {
	const Op = "smmproviderservice.FulfillOrder"

	// Get active providers from DB
	dbProviders, pgErr := s.repo.GetActiveByType(ctx, providerentity.ProviderTypeSMM)
	if pgErr != nil {
		return richerror.New(Op, pgErr)
	}

	// Find first provider with CLOSED circuit breaker
	var selectedProvider string
	var selectedAdapter SMMProvider
	for _, p := range dbProviders {
		adapter, ok := s.providers[p.Name]
		if !ok {
			continue
		}
		cb := s.breakers[p.Name]
		if cb.Allow() {
			selectedProvider = p.Name
			selectedAdapter = adapter
			break
		}
	}

	if selectedAdapter == nil {
		return fmt.Errorf("no available SMM provider")
	}

	// Create external order
	req := smmprams.CreateOrderAdapterRequest{
		ServiceID: fmt.Sprintf("%d", order.ProductID),
		Link:      order.TargetLink,
		Quantity:  order.Quantity,
	}
	resp, err := selectedAdapter.Create(ctx, req)
	if err != nil {
		s.breakers[selectedProvider].RecordFailure()
		return fmt.Errorf("provider %s failed: %w", selectedProvider, err)
	}

	s.breakers[selectedProvider].RecordSuccess()

	// Update order with external info
	// This would be done by the caller (orchestrator) via order service
	order.ExternalOrderID = resp.ExternalOrderID
	order.Status = orderentity.OrderStatusProcessing

	return nil
}
