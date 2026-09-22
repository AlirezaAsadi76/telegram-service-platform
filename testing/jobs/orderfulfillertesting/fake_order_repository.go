package orderfulfillertesting

import (
	"context"
	"errors"
	"time"

	"telegram-service-platform/entity/orderentity"
)

type fakeOrderRepository struct {
	order             *orderentity.Order
	claimResult       bool
	claimErr          error
	saveExternalErr   error
	createAttemptErr  error
	resolveAttemptErr error

	events []string

	attempts []*orderentity.FulfillmentAttempt
}

func (f *fakeOrderRepository) Create(_ context.Context, _ *orderentity.Order) error {
	return nil
}

func (f *fakeOrderRepository) GetByID(_ context.Context, orderID uint64) (*orderentity.Order, error) {
	f.events = append(f.events, "get_order")

	if f.order == nil || f.order.ID != orderID {
		return nil, errors.New("order not found")
	}

	return f.order, nil
}

func (f *fakeOrderRepository) GetStaleProcessingWithoutEvidence(
	_ context.Context,
	_ time.Duration,
	_ int,
) ([]*orderentity.Order, error) {
	return nil, nil
}

func (f *fakeOrderRepository) GetByUserID(_ context.Context, _ uint64) ([]*orderentity.Order, error) {
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

func (f *fakeOrderRepository) GetByStatus(_ context.Context, _ orderentity.OrderStatus) ([]*orderentity.Order, error) {
	return nil, nil
}

func (f *fakeOrderRepository) GetStalePaid(_ context.Context, _ time.Duration, _ int) ([]*orderentity.Order, error) {
	return nil, nil
}

func (f *fakeOrderRepository) ClaimForProcessing(_ context.Context, orderID uint64) (bool, error) {
	f.events = append(f.events, "claim")

	if f.claimErr != nil {
		return false, f.claimErr
	}

	if f.order == nil || f.order.ID != orderID {
		return false, errors.New("order not found")
	}

	if !f.claimResult {
		return false, nil
	}

	f.order.Status = orderentity.OrderStatusProcessing

	return true, nil
}

func (f *fakeOrderRepository) SaveExternalOrder(_ context.Context, orderID uint64, providerID uint64, externalOrderID string) error {
	f.events = append(f.events, "save_external_order")

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

func (f *fakeOrderRepository) AssignProvider(_ context.Context, orderID uint64, providerID uint64) error {
	f.events = append(f.events, "assign_provider")

	if f.order == nil || f.order.ID != orderID {
		return errors.New("order not found")
	}

	f.order.ProviderID = &providerID

	return nil
}

func (f *fakeOrderRepository) CreateFulfillmentAttempt(_ context.Context, attempt *orderentity.FulfillmentAttempt) (uint64, error) {
	f.events = append(f.events, "create_attempt")

	if f.createAttemptErr != nil {
		return 0, f.createAttemptErr
	}

	attempt.ID = uint64(len(f.attempts) + 1)

	f.attempts = append(
		f.attempts,
		attempt,
	)

	return attempt.ID, nil
}

func (f *fakeOrderRepository) GetUnresolvedFulfillmentAttempts(_ context.Context, _ int) ([]orderentity.FulfillmentAttempt, error) {
	return nil, nil
}

func (f *fakeOrderRepository) MarkFulfillmentAttemptResolved(_ context.Context, attemptID uint64) error {
	f.events = append(f.events, "resolve_attempt")

	if f.resolveAttemptErr != nil {
		return f.resolveAttemptErr
	}

	for _, attempt := range f.attempts {
		if attempt.ID == attemptID {
			now := time.Now()

			attempt.ResolvedAt = &now

			return nil
		}
	}

	return errors.New("attempt not found")
}

func (f *fakeOrderRepository) CompleteProcessing(_ context.Context, _ uint64) (bool, error) {
	return false, nil
}
