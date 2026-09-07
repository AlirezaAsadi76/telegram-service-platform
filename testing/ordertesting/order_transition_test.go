package orderentity

import (
	"errors"
	"telegram-service-platform/entity/orderentity"
	"testing"
)

func TestOrderStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name string
		from orderentity.OrderStatus
		to   orderentity.OrderStatus
		want bool
	}{
		{
			name: "pending to paid",
			from: orderentity.OrderStatusPending,
			to:   orderentity.OrderStatusPaid,
			want: true,
		},
		{
			name: "pending to expired",
			from: orderentity.OrderStatusPending,
			to:   orderentity.OrderStatusExpired,
			want: true,
		},
		{
			name: "pending to canceled",
			from: orderentity.OrderStatusPending,
			to:   orderentity.OrderStatusCanceled,
			want: true,
		},
		{
			name: "paid to processing",
			from: orderentity.OrderStatusPaid,
			to:   orderentity.OrderStatusProcessing,
			want: true,
		},
		{
			name: "processing to completed",
			from: orderentity.OrderStatusProcessing,
			to:   orderentity.OrderStatusCompleted,
			want: true,
		},
		{
			name: "processing to failed",
			from: orderentity.OrderStatusProcessing,
			to:   orderentity.OrderStatusFailed,
			want: true,
		},
		{
			name: "completed to paid",
			from: orderentity.OrderStatusCompleted,
			to:   orderentity.OrderStatusPaid,
			want: false,
		},
		{
			name: "failed to processing",
			from: orderentity.OrderStatusFailed,
			to:   orderentity.OrderStatusProcessing,
			want: false,
		},
		{
			name: "expired to paid",
			from: orderentity.OrderStatusExpired,
			to:   orderentity.OrderStatusPaid,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.from.CanTransitionTo(tt.to)
			if got != tt.want {
				t.Errorf("OrderStatus.CanTransitionTo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderStatus_ValidateTransitionTo(t *testing.T) {
	err := orderentity.OrderStatusCompleted.ValidateTransitionTo(orderentity.OrderStatusPaid)

	if !errors.Is(err, orderentity.ErrInvalidOrderTransition) {
		t.Fatal("expected ErrInvalidOrderTransition")
	}
}
