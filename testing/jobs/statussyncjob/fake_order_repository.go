package statussynctesting

import (
	"context"
	"errors"
	"time"

	"telegram-service-platform/entity/orderentity"
)

type fakeOrderRepository struct {
	orders []*orderentity.Order

	requestedStatus orderentity.OrderStatus

	completeResult bool
	completeErr    error

	events []string
}

func (f *fakeOrderRepository) Create(
	_ context.Context,
	_ *orderentity.Order,
) error {
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
	for _, order := range f.orders {
		if order.ID == orderID {
			return order, nil
		}
	}

	return nil, errors.New("order not found")
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
	status orderentity.OrderStatus,
) ([]*orderentity.Order, error) {
	f.events = append(f.events, "get_by_status")

	f.requestedStatus = status

	return f.orders, nil
}

func (f *fakeOrderRepository) GetStalePaid(
	_ context.Context,
	_ time.Duration,
	_ int,
) ([]*orderentity.Order, error) {
	return nil, nil
}

func (f *fakeOrderRepository) ClaimForProcessing(
	_ context.Context,
	_ uint64,
) (bool, error) {
	return false, nil
}

func (f *fakeOrderRepository) SaveExternalOrder(
	_ context.Context,
	_ uint64,
	_ uint64,
	_ string,
) error {
	return nil
}

func (f *fakeOrderRepository) AssignProvider(
	_ context.Context,
	_ uint64,
	_ uint64,
) error {
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
	return nil, nil
}

func (f *fakeOrderRepository) MarkFulfillmentAttemptResolved(
	_ context.Context,
	_ uint64,
) error {
	return nil
}

func (f *fakeOrderRepository) CompleteProcessing(
	_ context.Context,
	orderID uint64,
) (bool, error) {
	f.events = append(f.events, "complete_processing")

	if f.completeErr != nil {
		return false, f.completeErr
	}

	for _, order := range f.orders {
		if order.ID == orderID {
			if !f.completeResult {
				return false, nil
			}

			order.Status = orderentity.OrderStatusCompleted

			return true, nil
		}
	}

	return false, errors.New("order not found")
}
