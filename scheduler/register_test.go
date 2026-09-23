package scheduler

import (
	"testing"
	"time"
)

func TestScheduler_ResolveInterval_OrderFulfillmentRecovery(t *testing.T) {
	expected := 7 * time.Minute

	s := &Scheduler{
		config: Config{
			OrderFulfillmentRecoveryInterval: expected,
		},
	}

	got := s.resolveInterval(
		"order-fulfillment-recovery",
	)

	if got != expected {
		t.Fatalf(
			"expected recovery interval %s, got %s",
			expected,
			got,
		)
	}
}
