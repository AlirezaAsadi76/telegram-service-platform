package paymententity

import "errors"

var ErrInvalidPaymentTransition = errors.New("invalid payment status transition")

func (s PaymentStatus) CanTransitionTo(next PaymentStatus) bool {
	switch s {
	case PaymentStatusCreating:
		return next == PaymentStatusPending ||
			next == PaymentStatusUnknown ||
			next == PaymentStatusFailed

	case PaymentStatusPending:
		return next == PaymentStatusProcessing ||
			next == PaymentStatusSuccess ||
			next == PaymentStatusFailed ||
			next == PaymentStatusCanceled ||
			next == PaymentStatusExpired ||
			next == PaymentStatusUnknown

	case PaymentStatusProcessing:
		return next == PaymentStatusSuccess ||
			next == PaymentStatusFailed ||
			next == PaymentStatusUnknown

	case PaymentStatusUnknown:
		return next == PaymentStatusPending ||
			next == PaymentStatusProcessing ||
			next == PaymentStatusSuccess ||
			next == PaymentStatusFailed ||
			next == PaymentStatusCanceled ||
			next == PaymentStatusExpired

	case PaymentStatusSuccess,
		PaymentStatusFailed,
		PaymentStatusCanceled,
		PaymentStatusExpired:
		return false

	default:
		return false
	}
}

func (s PaymentStatus) ValidateTransitionTo(next PaymentStatus) error {
	if s.CanTransitionTo(next) {
		return nil
	}

	return ErrInvalidPaymentTransition
}
