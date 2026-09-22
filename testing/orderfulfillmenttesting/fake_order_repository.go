package orderfulfillmenttesting

import (
	"context"
	"errors"
	"telegram-service-platform/entity/orderentity"
	"time"
)

type fakeOrderRepository struct {
	orders map[uint64]*orderentity.Order

	attempts []*orderentity.FulfillmentAttempt

	events []string

	claimCalls              int
	saveExternalCalls       int
	createAttemptCalls      int
	resolveAttemptCalls     int
	assignProviderCalls     int
	completeProcessingCalls int
}

func newFakeOrderRepository(order *orderentity.Order) *fakeOrderRepository {
	return &fakeOrderRepository{
		orders: map[uint64]*orderentity.Order{
			order.ID: order,
		},
	}
}

func (f *fakeOrderRepository) Create(
	_ context.Context,
	order *orderentity.Order,
) error {
	f.orders[order.ID] = order

	return nil
}

func (f *fakeOrderRepository) GetStaleProcessingWithoutEvidence(
	_ context.Context,
	_ time.Duration,
	_ int,
) ([]*orderentity.Order, error) {
	return nil, nil
}

func (f *fakeOrderRepository) GetByID(
	_ context.Context,
	orderID uint64,
) (*orderentity.Order, error) {
	f.events = append(f.events, "get_order")

	order, ok := f.orders[orderID]
	if !ok {
		return nil, errors.New("order not found")
	}

	return order, nil
}

func (f *fakeOrderRepository) GetByUserID(
	_ context.Context,
	userID uint64,
) ([]*orderentity.Order, error) {
	result := make([]*orderentity.Order, 0)

	for _, order := range f.orders {
		if order.UserID == userID {
			result = append(result, order)
		}
	}

	return result, nil
}

func (f *fakeOrderRepository) UpdateStatus(
	_ context.Context,
	orderID uint64,
	status orderentity.OrderStatus,
	externalOrderID string,
	providerID *uint64,
) error {
	order, ok := f.orders[orderID]
	if !ok {
		return errors.New("order not found")
	}

	order.Status = status
	order.ExternalOrderID = externalOrderID
	order.ProviderID = providerID

	return nil
}

func (f *fakeOrderRepository) GetByStatus(
	_ context.Context,
	status orderentity.OrderStatus,
) ([]*orderentity.Order, error) {
	f.events = append(f.events, "get_by_status")

	result := make([]*orderentity.Order, 0)

	for _, order := range f.orders {
		if order.Status == status {
			result = append(result, order)
		}
	}

	return result, nil
}

func (f *fakeOrderRepository) GetStalePaid(
	_ context.Context,
	olderThan time.Duration,
	limit int,
) ([]*orderentity.Order, error) {
	f.events = append(f.events, "get_stale_paid")

	cutoff := time.Now().Add(-olderThan)

	result := make([]*orderentity.Order, 0)

	for _, order := range f.orders {
		if order.Status != orderentity.OrderStatusPaid {
			continue
		}

		if order.CreatedAt.IsZero() ||
			order.CreatedAt.Before(cutoff) {
			result = append(result, order)
		}

		if limit > 0 && len(result) >= limit {
			break
		}
	}

	return result, nil
}

func (f *fakeOrderRepository) ClaimForProcessing(
	_ context.Context,
	orderID uint64,
) (bool, error) {
	f.events = append(f.events, "claim")
	f.claimCalls++

	order, ok := f.orders[orderID]
	if !ok {
		return false, errors.New("order not found")
	}

	if order.Status != orderentity.OrderStatusPaid {
		return false, nil
	}

	order.Status = orderentity.OrderStatusProcessing

	return true, nil
}

func (f *fakeOrderRepository) SaveExternalOrder(
	_ context.Context,
	orderID uint64,
	providerID uint64,
	externalOrderID string,
) error {
	f.events = append(f.events, "save_external_order")
	f.saveExternalCalls++

	order, ok := f.orders[orderID]
	if !ok {
		return errors.New("order not found")
	}

	if order.Status != orderentity.OrderStatusProcessing {
		return errors.New("order is not processing")
	}

	if order.ProviderID != nil &&
		*order.ProviderID != providerID {
		return errors.New("provider mismatch")
	}

	if order.ExternalOrderID != "" &&
		order.ExternalOrderID != externalOrderID {
		return errors.New("external order ID mismatch")
	}

	order.ProviderID = &providerID
	order.ExternalOrderID = externalOrderID

	return nil
}

func (f *fakeOrderRepository) AssignProvider(
	_ context.Context,
	orderID uint64,
	providerID uint64,
) error {
	f.events = append(f.events, "assign_provider")
	f.assignProviderCalls++

	order, ok := f.orders[orderID]
	if !ok {
		return errors.New("order not found")
	}

	if order.Status != orderentity.OrderStatusProcessing {
		return errors.New("order is not processing")
	}

	order.ProviderID = &providerID

	return nil
}

func (f *fakeOrderRepository) CreateFulfillmentAttempt(
	_ context.Context,
	attempt *orderentity.FulfillmentAttempt,
) (uint64, error) {
	f.events = append(f.events, "create_attempt")
	f.createAttemptCalls++

	attempt.ID = uint64(len(f.attempts) + 1)

	if attempt.CreatedAt.IsZero() {
		attempt.CreatedAt = time.Now()
	}

	if attempt.UpdatedAt.IsZero() {
		attempt.UpdatedAt = attempt.CreatedAt
	}

	f.attempts = append(
		f.attempts,
		attempt,
	)

	return attempt.ID, nil
}

func (f *fakeOrderRepository) GetUnresolvedFulfillmentAttempts(
	_ context.Context,
	limit int,
) ([]orderentity.FulfillmentAttempt, error) {
	f.events = append(f.events, "get_unresolved_attempts")

	result := make([]orderentity.FulfillmentAttempt, 0)

	for _, attempt := range f.attempts {
		if attempt.ResolvedAt != nil {
			continue
		}

		result = append(result, *attempt)

		if limit > 0 && len(result) >= limit {
			break
		}
	}

	return result, nil
}

func (f *fakeOrderRepository) MarkFulfillmentAttemptResolved(
	_ context.Context,
	attemptID uint64,
) error {
	f.events = append(f.events, "resolve_attempt")
	f.resolveAttemptCalls++

	for _, attempt := range f.attempts {
		if attempt.ID != attemptID {
			continue
		}

		now := time.Now()
		attempt.ResolvedAt = &now
		attempt.UpdatedAt = now

		return nil
	}

	return errors.New("attempt not found")
}

func (f *fakeOrderRepository) CompleteProcessing(
	_ context.Context,
	orderID uint64,
) (bool, error) {
	f.events = append(f.events, "complete_processing")
	f.completeProcessingCalls++

	order, ok := f.orders[orderID]
	if !ok {
		return false, errors.New("order not found")
	}

	if order.Status != orderentity.OrderStatusProcessing {
		return false, nil
	}

	order.Status = orderentity.OrderStatusCompleted

	return true, nil
}
