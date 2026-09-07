package orderentity

import (
	"errors"
)

var ErrInvalidOrderTransition = errors.New("invalid order status transition")

func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	switch s {
	case OrderStatusPending:
		return next == OrderStatusPaid ||
			next == OrderStatusCanceled ||
			next == OrderStatusExpired

	case OrderStatusPaid:
		return next == OrderStatusProcessing ||
			next == OrderStatusCanceled

	case OrderStatusProcessing:
		return next == OrderStatusCompleted ||
			next == OrderStatusFailed

	case OrderStatusCompleted,
		OrderStatusFailed,
		OrderStatusCanceled,
		OrderStatusExpired:
		return false

	default:
		return false
	}
}

func (s OrderStatus) ValidateTransitionTo(next OrderStatus) error {
	if s.CanTransitionTo(next) {
		return nil
	}

	return ErrInvalidOrderTransition
}
