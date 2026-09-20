package orderfulfillmentrecoverytesting

import (
	"context"
	"errors"
	"time"

	"telegram-service-platform/entity/orderentity"
)

type fakeOrderRepository struct {
	order *orderentity.Order

	unresolvedAttempts []orderentity.FulfillmentAttempt

	saveExternalCalls   int
	assignProviderCalls int
	resolveCalls        int

	events []string

	saveExternalErr   error
	assignProviderErr error
	resolveErr        error

	staleOrders []*orderentity.Order
}

func (f *fakeOrderRepository) Create(
	_ context.Context,
	_ *orderentity.Order,
) error {
	return nil
}

func (f *fakeOrderRepository) GetByID(
	_ context.Context,
	orderID uint64,
) (*orderentity.Order, error) {
	f.events = append(f.events, "get_order")

	if f.order == nil || f.order.ID != orderID {
		return nil, errors.New("order not found")
	}

	return f.order, nil
}

func (f *fakeOrderRepository) GetByUserID(
	_ context.Context,
	_ uint64,
) ([]*orderentity.Order, error) {
	return nil, nil
}

func (f *fakeOrderRepository) UpdateStatus(
	_ context.Context,
	_ uint64,
	_ orderentity.OrderStatus,
	_ string,
	_ *uint64,
) error {
	return nil
}

func (f *fakeOrderRepository) GetByStatus(
	_ context.Context,
	_ orderentity.OrderStatus,
) ([]*orderentity.Order, error) {
	return nil, nil
}

func (f *fakeOrderRepository) GetStalePaid(
	_ context.Context,
	_ time.Duration,
	_ int,
) ([]*orderentity.Order, error) {
	f.events = append(f.events, "get_stale_paid")

	return f.staleOrders, nil
}

func (f *fakeOrderRepository) ClaimForProcessing(
	_ context.Context,
	_ uint64,
) (bool, error) {
	return false, nil
}

func (f *fakeOrderRepository) SaveExternalOrder(
	_ context.Context,
	orderID uint64,
	providerID uint64,
	externalOrderID string,
) error {
	f.events = append(f.events, "save_external_order")
	f.saveExternalCalls++

	if f.saveExternalErr != nil {
		return f.saveExternalErr
	}

	if f.order == nil || f.order.ID != orderID {
		return errors.New("order not found")
	}

	f.order.ProviderID = &providerID
	f.order.ExternalOrderID = externalOrderID

	return nil
}

func (f *fakeOrderRepository) AssignProvider(
	_ context.Context,
	orderID uint64,
	providerID uint64,
) error {
	f.events = append(f.events, "assign_provider")
	f.assignProviderCalls++

	if f.assignProviderErr != nil {
		return f.assignProviderErr
	}

	if f.order == nil || f.order.ID != orderID {
		return errors.New("order not found")
	}

	f.order.ProviderID = &providerID

	return nil
}

func (f *fakeOrderRepository) CreateFulfillmentAttempt(
	_ context.Context,
	_ *orderentity.FulfillmentAttempt,
) (uint64, error) {
	return 0, nil
}

func (f *fakeOrderRepository) GetUnresolvedFulfillmentAttempts(
	_ context.Context,
	_ int,
) ([]orderentity.FulfillmentAttempt, error) {
	f.events = append(f.events, "get_unresolved_attempts")

	return f.unresolvedAttempts, nil
}

func (f *fakeOrderRepository) MarkFulfillmentAttemptResolved(
	_ context.Context,
	attemptID uint64,
) error {
	f.events = append(f.events, "resolve_attempt")
	f.resolveCalls++

	if f.resolveErr != nil {
		return f.resolveErr
	}

	for i := range f.unresolvedAttempts {
		if f.unresolvedAttempts[i].ID == attemptID {
			now := time.Now()
			f.unresolvedAttempts[i].ResolvedAt = &now

			return nil
		}
	}

	return errors.New("attempt not found")
}

func (f *fakeOrderRepository) CompleteProcessing(
	_ context.Context,
	_ uint64,
) (bool, error) {
	return false, nil
}
