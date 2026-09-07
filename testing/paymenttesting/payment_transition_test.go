package paymenttesting

import (
	"errors"
	"telegram-service-platform/entity/paymententity"
	"testing"
)

func TestPaymentStatus_CanTransitionTo(t *testing.T) {

	tests := []struct {
		name string
		from paymententity.PaymentStatus
		to   paymententity.PaymentStatus
		want bool
	}{
		{
			name: "Creating to pending",
			from: paymententity.PaymentStatusCreating,
			to:   paymententity.PaymentStatusPending,
			want: true,
		},
		{
			name: "Creating to failed",
			from: paymententity.PaymentStatusCreating,
			to:   paymententity.PaymentStatusFailed,
			want: true,
		},
		{
			name: "Creating to unknown",
			from: paymententity.PaymentStatusCreating,
			to:   paymententity.PaymentStatusUnknown,
			want: true,
		},
		{
			name: "pending to expired",
			from: paymententity.PaymentStatusPending,
			to:   paymententity.PaymentStatusExpired,
			want: true,
		},
		{
			name: "Pending to unknown",
			from: paymententity.PaymentStatusPending,
			to:   paymententity.PaymentStatusUnknown,
			want: true,
		},
		{
			name: "Pending to processing",
			from: paymententity.PaymentStatusPending,
			to:   paymententity.PaymentStatusProcessing,
			want: true,
		},
		{
			name: "Pending to failed",
			from: paymententity.PaymentStatusPending,
			to:   paymententity.PaymentStatusFailed,
			want: true,
		},
		{
			name: "Pending to success",
			from: paymententity.PaymentStatusPending,
			to:   paymententity.PaymentStatusSuccess,
			want: true,
		},
		{
			name: "Pending to cancelled",
			from: paymententity.PaymentStatusPending,
			to:   paymententity.PaymentStatusCanceled,
			want: true,
		},
		{
			name: "processing to unknown",
			from: paymententity.PaymentStatusProcessing,
			to:   paymententity.PaymentStatusUnknown,
			want: true,
		},
		{
			name: "processing to failed",
			from: paymententity.PaymentStatusProcessing,
			to:   paymententity.PaymentStatusFailed,
			want: true,
		},
		{
			name: "processing to success",
			from: paymententity.PaymentStatusProcessing,
			to:   paymententity.PaymentStatusSuccess,
			want: true,
		},
		{
			name: "unknow to expired",
			from: paymententity.PaymentStatusUnknown,
			to:   paymententity.PaymentStatusExpired,
			want: true,
		},
		{
			name: "unknow to failed",
			from: paymententity.PaymentStatusUnknown,
			to:   paymententity.PaymentStatusFailed,
			want: true,
		},
		{
			name: "unknow to cancelled",
			from: paymententity.PaymentStatusUnknown,
			to:   paymententity.PaymentStatusCanceled,
			want: true,
		},
		{
			name: "unknow to pending",
			from: paymententity.PaymentStatusUnknown,
			to:   paymententity.PaymentStatusPending,
			want: true,
		},
		{
			name: "unknow to success",
			from: paymententity.PaymentStatusUnknown,
			to:   paymententity.PaymentStatusSuccess,
			want: true,
		},
		{
			name: "unknow to pending",
			from: paymententity.PaymentStatusUnknown,
			to:   paymententity.PaymentStatusPending,
			want: true,
		},
		{
			name: "success to pending",
			from: paymententity.PaymentStatusSuccess,
			to:   paymententity.PaymentStatusPending,
			want: false,
		},
		{
			name: "success to failed",
			from: paymententity.PaymentStatusSuccess,
			to:   paymententity.PaymentStatusFailed,
			want: false,
		},
		{
			name: "failed to cancelled",
			from: paymententity.PaymentStatusFailed,
			to:   paymententity.PaymentStatusCanceled,
			want: false,
		},
		{
			name: "canceled to pending",
			from: paymententity.PaymentStatusCanceled,
			to:   paymententity.PaymentStatusPending,
			want: false,
		},
		{
			name: "expired to success",
			from: paymententity.PaymentStatusExpired,
			to:   paymententity.PaymentStatusSuccess,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
				t.Errorf("CanTransitionTo() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestPaymentStatus_ValidateTransitionTo(t *testing.T) {
	err := paymententity.PaymentStatusCreating.ValidateTransitionTo(paymententity.PaymentStatusProcessing)

	if !errors.Is(err, paymententity.ErrInvalidPaymentTransition) {
		t.Fatal("expected ErrInvalidPaymentTransition")
	}
}
